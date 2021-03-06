////////////////////////////////////////////////////////////////////////////////
// Manages various aspects e.g. insert, update, reset of the master key table //
////////////////////////////////////////////////////////////////////////////////

package masterpasswordmanager

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dasdipanjan04/gpwm/connect"
	"github.com/dasdipanjan04/gpwm/gpwmcrypto"
	"github.com/dasdipanjan04/gpwm/helper/glogger"
	"github.com/dasdipanjan04/gpwm/helper/gqrpdf"
	"github.com/dasdipanjan04/gpwm/helper/gretry"
	"github.com/dasdipanjan04/gpwm/helper/gscan"

	_ "github.com/lib/pq" // imports lib pq for psql
)

type userDetails struct {
	firstName    string
	lastName     string
	email        string
	password     string
	oldMasterkey string
}

// CreateMasterKeyTable creates a mastertable if it doesn't exist.
func CreateMasterKeyTable() *sql.DB {
	createTable := `create table if not exists mastertable (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		master_key BYTEA NOT NULL,
		created_at TEXT,
		updated_at TEXT,
		is_active BOOL
	  );`

	db, err := connect.OpenDB()
	if err != nil {
		glogger.Glog("masterkeymanager:CreateMasterKeyTable:OpenDB ", err.Error())
		return nil
	}

	_, err = db.Exec(createTable)
	if err != nil {
		glogger.Glog("masterkeymanager:CreateMasterKeyTable:Exec ", err.Error())
		return nil
	}

	return db
}

// InsertUserDataToDB inserts new information to the database.
func InsertUserDataToDB(db *sql.DB, firstName string, lastName string,
	email string, isActive bool) {

	insertStatement := `INSERT INTO mastertable (first_name, last_name, email, master_key, created_at, updated_at, is_active)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT email FROM mastertable where(mastertable.email = $3));`

	//Ask user to set a Master Password
	fmt.Println("Please Set a Strong Master Password:")
	masterPassword := gscan.GscanFromTerminal()

	// Randomly generate master key for each user data insert to the masterkey table.
	masterAccountKey := gpwmcrypto.GenerateAccountSecretKey()
	masterAccountKeyByte := []byte(masterAccountKey)
	encryptedMasterAccountKey, eerr := gpwmcrypto.EncryptKEKAES(masterAccountKeyByte, masterPassword, email)
	if eerr != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:EncryptMasterKEKAES ", eerr.Error())
	}

	timeNow := time.Now().Unix()
	createdAt := strconv.FormatInt(timeNow, 10)
	updatedAt := strconv.FormatInt(timeNow, 10)

	_, err := db.Exec(insertStatement, firstName, lastName, email, encryptedMasterAccountKey, createdAt, updatedAt, isActive)
	if err != nil {
		glogger.Glog("masterkeymanager:InsertMasterKeyDataToDB:Exec ", err.Error())
		return
	}

	gqrpdf.MasterKeyQRCodePDFGenerator(masterAccountKey, firstName, lastName)
}

// UpdateInfo updates the user database info if changed.
func UpdateInfo(db *sql.DB, id int, firstName string, lastName string,
	email string, masterKey string, createdAt string,
	updatedAt string, isActive bool) {

	updateStatement := `UPDATE mastertable 
	SET first_name = $2, last_name = $3, email = $4, master_key = $5, created_at = $6, updated_at = $7, is_active = $8
	WHERE id = $1;`

	_, err := db.Exec(updateStatement, id, firstName, lastName, email, masterKey, createdAt, updatedAt, isActive)
	if err != nil {
		glogger.Glog("masterkeymanager:UpdateInfo:Exec ", err.Error())
		return
	}
}

// ResetMasterKey resets master key in the database.
func ResetMasterKey(db *sql.DB) error {

	userdetail, err := GetUserDetails(db)

	gretry.MAXIMUMALLOWEDATTEMPTS--

	if err != nil {
		rerr := gretry.Retry(func(attempts int) error {
			reseterr := ResetMasterKey(db)
			return reseterr
		})
		if rerr != nil {
			log.Fatalln(rerr.Error())
			return rerr
		}
		return rerr
	}

	// Finding the id and the stored encrypted master key from the table.
	id := 0
	oldMasterKeyFromTable := []byte("")
	findIdByEmail := fmt.Sprintf(`SELECT id, master_key FROM mastertable WHERE email in (%s);`, userdetail.email)
	err = db.QueryRow(findIdByEmail).Scan(&id, &oldMasterKeyFromTable)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:QueryRow ", err.Error())
		return err
	}

	emailId := strings.Trim(userdetail.email, "'")

	// decrypt oldmasterkey and compare
	dycryptText, err := gpwmcrypto.DecryptAESKEK(oldMasterKeyFromTable, string(userdetail.password), emailId)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKEK ", err.Error())
		return err
	}

	if dycryptText != userdetail.oldMasterkey {
		err = errors.New("key doesn't match")
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKEK ", err.Error())
		return err
	}

	// Generate new account key or master key.
	newMasterAccountKey := gpwmcrypto.GenerateAccountSecretKey()
	// Encrypt new master key.
	encryptedText, err := gpwmcrypto.EncryptKEKAES([]byte(newMasterAccountKey), userdetail.password, emailId)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:EncryptMasterKEKAES ", err.Error())
		return err
	}

	// All the information is correct until this point, trying to reset the master key.
	resetMasterKeyStatement := `UPDATE mastertable
	SET master_key = $1
	WHERE id = $2;`
	_, err = db.Exec(resetMasterKeyStatement, encryptedText, id)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:Exec ", err.Error())
		return err
	}
	glogger.Glog("masterkeymanager:ResetMasterKey ", "You have successfully reset your master key")

	// Generate QR Code with new master key
	gqrpdf.MasterKeyQRCodePDFGenerator(newMasterAccountKey, userdetail.firstName, userdetail.lastName)
	return err
}

//GetUserDetails get the details of the user.
func GetUserDetails(db *sql.DB) (*userDetails, error) {

	fmt.Println("Reset your masterkey")
	fmt.Println("Please enter your registered email address:")
	email := gscan.GscanFromTerminal()
	email = "'" + email + "'"
	id := 0
	findIdByEmail := fmt.Sprintf(`SELECT id FROM mastertable WHERE email in (%s);`, email)
	err := db.QueryRow(findIdByEmail).Scan(&id)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:QueryRow ", err.Error())
		fmt.Println("Sorry! your email address is incorrect. Please try again with the correct email address")
		return nil, err
	}

	firstName := ""
	lastName := ""
	findNameById := fmt.Sprintf(`SELECT first_name, last_name FROM mastertable WHERE id in (%d)`, id)
	err = db.QueryRow(findNameById).Scan(&firstName, &lastName)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:QueryRow ", err.Error())
		fmt.Println("Unable to retrieve name from the database.")
		return nil, err
	}

	fmt.Println("Please enter your password:")
	password := gscan.GscanFromTerminal()

	fmt.Println("Please enter your current master key pass:")
	oldMasterKey := gscan.GscanFromTerminal()

	userdetail := userDetails{}
	userdetail.firstName = firstName
	userdetail.lastName = lastName
	userdetail.email = email
	userdetail.password = password
	userdetail.oldMasterkey = oldMasterKey

	return &userdetail, err
}

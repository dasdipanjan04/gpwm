////////////////////////////////////////////////////////////////////////////////
// Manages various aspects e.g. insert, update, reset of the master key table //
////////////////////////////////////////////////////////////////////////////////

package masterkeymanager

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dasdipanjan04/gpwm/connect"
	"github.com/dasdipanjan04/gpwm/helper/glogger"
	"github.com/dasdipanjan04/gpwm/helper/gqrpdf"
	"github.com/dasdipanjan04/gpwm/helper/gretry"
	"github.com/dasdipanjan04/gpwm/helper/gscan"
	"github.com/dasdipanjan04/gpwm/masterkeysecure"

	_ "github.com/lib/pq"
)

type userDetails struct {
	email        string
	password     string
	oldMasterkey string
	newMasterkey string
}

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

func InsertUserDataToDB(db *sql.DB, first_name string, last_name string,
	email string, password string, is_active bool) {

	insertStatement := `INSERT INTO mastertable (first_name, last_name, email, master_key, created_at, updated_at, is_active)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT email FROM mastertable where(mastertable.email = $3));`

	// Randomly generate master key for each user data insert to the masterkey table.
	master_key := GenerateMasterKey()
	master_key_byte := []byte(master_key)
	encrypted_master_key, eerr := masterkeysecure.EncryptMasterKeyAES(master_key_byte, password)
	if eerr != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:EncryptMasterKeyAES ", eerr.Error())
	}

	time_now := time.Now().Unix()
	created_at := strconv.FormatInt(time_now, 10)
	updated_at := strconv.FormatInt(time_now, 10)

	_, err := db.Exec(insertStatement, first_name, last_name, email, encrypted_master_key, created_at, updated_at, is_active)
	if err != nil {
		glogger.Glog("masterkeymanager:InsertMasterKeyDataToDB:Exec ", err.Error())
		return
	}

	gqrpdf.MasterKeyQRCodePDFGenerator(master_key, first_name, last_name)
}

func UpdateInfo(db *sql.DB, id int, first_name string, last_name string,
	email string, master_key string, created_at string,
	updated_at string, is_active bool) {

	updateStatement := `UPDATE mastertable 
	SET first_name = $2, last_name = $3, email = $4, master_key = $5, created_at = $6, updated_at = $7, is_active = $8
	WHERE id = $1;`

	_, err := db.Exec(updateStatement, id, first_name, last_name, email, master_key, created_at, updated_at, is_active)
	if err != nil {
		glogger.Glog("masterkeymanager:UpdateInfo:Exec ", err.Error())
		return
	}
}

// Resets master key in the database.
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

	// decrypt oldmasterkey and compare
	dycryptText, derr := masterkeysecure.DecryptAESMasterKey(oldMasterKeyFromTable, userdetail.password)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKey ", derr.Error())
		return err
	}
	if dycryptText != userdetail.oldMasterkey {
		err = errors.New("Key doesn't match")
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKey ", err.Error())
		return err
	}

	// encrypt new master key.
	encryptedText, err := masterkeysecure.EncryptMasterKeyAES([]byte(userdetail.newMasterkey), userdetail.password)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:EncryptMasterKeyAES ", err.Error())
		return err
	}

	// All the information is correct until this point, trying to reset the master key.
	reserMasterKeyStatement := `UPDATE mastertable
	SET master_key = $1
	WHERE id = $2;`
	_, err = db.Exec(reserMasterKeyStatement, encryptedText, id)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:Exec ", err.Error())
		return err
	}
	glogger.Glog("masterkeymanager:ResetMasterKey ", "You have successfully reset your master key")

	return err
}

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

	fmt.Println("Please enter your password:")
	password := gscan.GscanFromTerminal()

	fmt.Println("Please enter your current master key pass:")
	oldMasterKey := gscan.GscanFromTerminal()

	fmt.Println("Please enter your new master key pass:")
	newMasterKey := gscan.GscanFromTerminal()

	userdetail := userDetails{}
	userdetail.email = email
	userdetail.password = password
	userdetail.oldMasterkey = oldMasterKey
	userdetail.newMasterkey = newMasterKey

	return &userdetail, err
}

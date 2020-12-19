////////////////////////////////////////////////////////////////////////////////
// Manages various aspects e.g. insert, update, reset of the master key table //
////////////////////////////////////////////////////////////////////////////////

package masterkeymanager

import (
	"database/sql"
	"fmt"
	"gpwm/connect"
	"gpwm/internal/glogger"
	"gpwm/masterkeysecure"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

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

func InsertMasterKeyDataToDB(db *sql.DB, first_name string, last_name string,
	email string, master_key string, password string, is_active bool) {

	insertStatement := `INSERT INTO mastertable (first_name, last_name, email, master_key, created_at, updated_at, is_active)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT email FROM mastertable where(mastertable.email = $3));`

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
func ResetMasterKey(db *sql.DB, email string, password string, oldMasterKey string, newmasterKey string, masterKeylost bool) {

	findIdByEmail := fmt.Sprintf(`SELECT id, master_key FROM mastertable WHERE email in (%s);`, email)

	id := 0
	oldMasterKeyFromTable := []byte("")
	err := db.QueryRow(findIdByEmail).Scan(&id, &oldMasterKeyFromTable)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:QueryRow ", err.Error())
		return
	}
	// decrypt oldmasterkey and compare
	dycryptText, derr := masterkeysecure.DecryptAESMasterKey(oldMasterKeyFromTable, password)
	if dycryptText != oldMasterKey || derr != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKey ", derr.Error())
		return
	}

	// encrypt new master key.
	encryptedText, eerr := masterkeysecure.EncryptMasterKeyAES([]byte(newmasterKey), password)
	if eerr != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:EncryptMasterKeyAES ", eerr.Error())
	}

	reserMasterKeyStatement := `UPDATE mastertable
	SET master_key = $1
	WHERE id = $2;`

	_, err = db.Exec(reserMasterKeyStatement, encryptedText, id)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:Exec ", err.Error())
		return
	}

	glogger.Glog("masterkeymanager:ResetMasterKey ", "You have successfully reset your master key")
}

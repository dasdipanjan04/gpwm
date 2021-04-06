package passwordmanager

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dasdipanjan04/gpwm/connect"
	"github.com/dasdipanjan04/gpwm/gpwmcrypto"
	"github.com/dasdipanjan04/gpwm/helper/glogger"
)

// CreatePasswordManagerTable creates a password manager table if it doesn't exist.
func CreatePasswordManagerTable() *sql.DB {
	createTable := `create table if not exists passwordmanagertable (
		id SERIAL PRIMARY KEY,
		application_name TEXT UNIQUE NOT NULL,
		user_name TEXT UNIQUE NOT NULL,
		apllication_password BYTEA NOT NULL,
		created_at TEXT,
		updated_at TEXT,
		is_active BOOL
	  );`

	passwordmanagerDb, err := connect.OpenDB()
	if err != nil {
		glogger.Glog("passwordmanager:CreatePasswordManagerTable:OpenDB ", err.Error())
		return nil
	}

	_, err = passwordmanagerDb.Exec(createTable)
	if err != nil {
		glogger.Glog("passwordmanager:CreatePasswordManagerTable:Exec ", err.Error())
		return nil
	}

	return passwordmanagerDb
}

// InsertEncryptedPasswordToDB inserts newly encrypted app password with user name to the database.
func InsertEncryptedPasswordToDB(dbPassmanager *sql.DB,
	applicationName string,
	username string,
	encryptedAppPassword []byte,
	isActive bool) error {

	insertStatement := `INSERT INTO passwordmanagertable (application_name, user_name, apllication_password, created_at, updated_at, is_active)
	SELECT $1, $2, $3, $4, $5, $6
	WHERE NOT EXISTS (SELECT application_name FROM passwordmanagertable where(passwordmanagertable.application_name = $1));`
	_, err := dbPassmanager.Exec(insertStatement, applicationName, "1test@test.test", encryptedAppPassword, "created_at", "updated_at", true)
	if err != nil {
		glogger.Glog("passwordmanager:InsertMasterKeyDataToDB:Exec ", err.Error())
		return err
	}
	return err
}

// EncryptApplicationPassword encrypts application password.
func EncryptApplicationPassword(dbMaster *sql.DB,
	dbPassmanager *sql.DB,
	masterPassword string,
	email string,
	application string,
	appPassword string) ([]byte, error) {

	id := 0
	oldMasterKeyFromTable := []byte("")
	findIdByEmail := fmt.Sprintf(`SELECT id, master_key FROM mastertable WHERE email in (%s);`, email)
	err := dbMaster.QueryRow(findIdByEmail).Scan(&id, &oldMasterKeyFromTable)
	if err != nil {
		glogger.Glog("passwordmanager:EncryptApplicationPassword:QueryRow ", err.Error())
		return nil, err
	}

	email = strings.Trim(email, "'")
	// decrypt oldmasterkey and compare
	dycryptedMasterKey, err := gpwmcrypto.DecryptAESKEK(oldMasterKeyFromTable, masterPassword, email)
	if err != nil {
		glogger.Glog("passwordmanager:EncryptApplicationPassword:DecryptAESMasterKEK ", err.Error())
		return nil, err
	}

	applicationPasswordByte := []byte(appPassword)
	encryptedAppPassword, err := gpwmcrypto.EncryptKEKAES(applicationPasswordByte, dycryptedMasterKey, email)
	if err != nil {
		glogger.Glog("passwordmanager:EncryptApplicationPassword:EncryptMasterKEKAES ", err.Error())
		return nil, err
	}

	return encryptedAppPassword, err
}

// DecryptAppPassword decrypts the application password.
func DecryptAppPassword(encryptedAppPassword []byte, masterKey string, email string) (string, error) {
	email = strings.Trim(email, "'")
	dycryptAppPass, err := gpwmcrypto.DecryptAESKEK(encryptedAppPassword, masterKey, email)
	if err != nil {
		glogger.Glog("passwordmanager:DecryptAppPassword:DecryptAESMasterKEK ", err.Error())
		return "", err
	}
	fmt.Println("Dencrypted App Password: ", dycryptAppPass)
	return dycryptAppPass, err
}

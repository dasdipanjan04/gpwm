package passwordmanager

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dasdipanjan04/gpwm/helper/glogger"
	"github.com/dasdipanjan04/gpwm/masterkeysecure"
)

func EncryptApplicationPassword(db *sql.DB, masterPassword string, email string, application string, appPassword string) error {
	id := 0
	oldMasterKeyFromTable := []byte("")
	findIdByEmail := fmt.Sprintf(`SELECT id, master_key FROM mastertable WHERE email in (%s);`, email)
	err := db.QueryRow(findIdByEmail).Scan(&id, &oldMasterKeyFromTable)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:QueryRow ", err.Error())
		return err
	}

	email = strings.Trim(email, "'")
	// decrypt oldmasterkey and compare
	dycryptedMasterKey, err := masterkeysecure.DecryptAESMasterKEK(oldMasterKeyFromTable, masterPassword, email)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKEK ", err.Error())
		return err
	}

	application_password_byte := []byte(appPassword)
	encrypted_app_password, err := masterkeysecure.EncryptMasterKEKAES(application_password_byte, dycryptedMasterKey, email)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:EncryptMasterKEKAES ", err.Error())
		return err
	}

	dycryptAppPass, err := masterkeysecure.DecryptAESMasterKEK(encrypted_app_password, dycryptedMasterKey, email)
	if err != nil {
		glogger.Glog("masterkeymanager:ResetMasterKey:DecryptAESMasterKEK ", err.Error())
		return err
	}
	fmt.Println("Dencrypted App Password: ", dycryptAppPass)

	return err
}

package main

import (
	"github.com/dasdipanjan04/gpwm/helper/glogger"
	mkm "github.com/dasdipanjan04/gpwm/masterpasswordmanager"
	"github.com/dasdipanjan04/gpwm/passwordmanager"
)

func main() {
	db := mkm.CreateMasterKeyTable()
	mkm.InsertUserDataToDB(db, "TestFN1", "TestLN1", "1test@test.test", true)
	mkm.ResetMasterKey(db)

	passwordManagerTable := passwordmanager.CreatePasswordManagerTable()
	encrypted_app_password, err := passwordmanager.EncryptApplicationPassword(db, passwordManagerTable,
		"1234",
		"'1test@test.test'",
		"facebook",
		"fbpass")
	if err != nil {
		glogger.Glog("Main:EncryptApplicationPassword", err.Error())
	}
	passwordmanager.InsertEncryptedPasswordToDB(passwordManagerTable,
		"facebook",
		"'1test@test.test'",
		encrypted_app_password,
		true)
}

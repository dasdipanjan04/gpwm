package main

import (
	mkm "github.com/dasdipanjan04/gpwm/masterpasswordmanager"
)

func main() {
	db := mkm.CreateMasterKeyTable()
	mkm.InsertUserDataToDB(db, "TestFN_1", "TestLN_1", "1test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_2", "TestLN_2", "2test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_3", "TestLN_3", "3test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_4", "TestLN_4", "4test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_5", "TestLN_5", "5test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_6", "TestLN_6", "6test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_7", "TestLN_7", "7test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_8", "TestLN_8", "8test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_9", "TestLN_9", "9test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_10", "TestLN_10", "10test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_11", "TestLN_11", "11test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_12", "TestLN_12", "12test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_13", "TestLN_13", "13test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_14", "TestLN_14", "14test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_15", "TestLN_15", "15test@test.test", true)
	mkm.InsertUserDataToDB(db, "TestFN_16", "TestLN_16", "16test@test.test", true)
	mkm.ResetMasterKey(db)
}

package main

import (
	mkm "gpwm/masterkeymanager"
	"strconv"
	"time"
)

func main() {
	db := mkm.CreateMasterKeyTable()
	err := mkm.InsertMasterKeyDataToDB(db, "TestFN", "TestLN", "test@test.test", "Test", true)
	if err != nil {
		panic(err)
	}

	err = mkm.InsertMasterKeyDataToDB(db, "TestFN_2", "TestLN_2", "testtest@testtest.testtest", "TestTest", true)
	if err != nil {
		panic(err)
	}

	time_now := time.Now().Unix()
	created_at := strconv.FormatInt(time_now, 10)
	updated_at := strconv.FormatInt(time_now, 10)
	err = mkm.UpdateInfo(db, 1, "TesFN", "TestLN", "test@test.test", "FooBar", created_at, updated_at, true)
	if err != nil {
		panic(err)
	}

	emailstr1 := "'test@tsadsdfsdafasdfest.test'"
	emailstr2 := "'test@test.test'"

	mkm.ResetMasterKey(db, emailstr1, "'testtestetststststststststesttestetstststststststs'")
	mkm.ResetMasterKey(db, emailstr2, "'TESSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSTTTTTTTTTTTTTTTTT'")
}

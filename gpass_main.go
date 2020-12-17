package main

import (
	"gpwm/connect"
	"strconv"
	"time"
)

func main() {
	db := connect.CreateMasterKeyTable()
	err := connect.InsertMasterKeyDataToDB(db, "TestFN", "TestLN", "test@test.test", "Test", true)
	if err != nil {
		panic(err)
	}

	err = connect.InsertMasterKeyDataToDB(db, "TestFN_2", "TestLN_2", "testtest@testtest.testtest", "TestTest", true)
	if err != nil {
		panic(err)
	}

	time_now := time.Now().Unix()
	created_at := strconv.FormatInt(time_now, 10)
	updated_at := strconv.FormatInt(time_now, 10)
	err = connect.UpdateInfo(db, 1, "TesFN", "TestLN", "test@test.test", "FooBar", created_at, updated_at, true)
	if err != nil {
		panic(err)
	}

	emailstr1 := "'test@tsadsdfsdafasdfest.test'"
	emailstr2 := "'test@test.test'"

	connect.ResetMasterKey(db, emailstr1, "'testtestetststststststststesttestetstststststststs'")
	connect.ResetMasterKey(db, emailstr2, "'FooBarFooBarFooBarFooBarFooBarFooBarFooBarFooBar'")
}

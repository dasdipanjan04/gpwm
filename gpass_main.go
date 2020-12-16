package main

import (
	"gpwm/connect"
	"strconv"
	"time"
)

func main() {
	db := connect.CreateMasterKeyTable()
	err := connect.InsertMasterKeyDataToDB(db, "TesFNffgsfg", "TestLN", "firsttfsgsfgsfdest@test.test", "Test", true)
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
}

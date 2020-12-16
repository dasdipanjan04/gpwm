package main

import (
	"fmt"
	"gpwm/connect"
)

func main() {
	db := connect.CreateMasterKeyTable()
	err := connect.InsertMasterKeyDataToDB(db, "TesFN", "TestLN", "firsttest@test.test", "Test", true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table Created")
}

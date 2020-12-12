package main

import (
	"fmt"

	"gpwm/connect"
)

func main() {
	fmt.Println("ok")
	connect.ConnectToMasterDB()
}

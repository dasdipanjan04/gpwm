package main

import (
	"fmt"
	"gpwm/connect"
)

func main() {
	connect.CreateTable()
	fmt.Println("Table Created")
}

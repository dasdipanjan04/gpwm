package main

import (
	"gpwm/connect"

	_ "github.com/lib/pq"
)

func CreateTable() {
	createTable := `create table if not exists mastertable (
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		master_key TEXT NOT NULL,
		created_at TEXT,
		updated_at TEXT,
		is_active BOOL
	  );`

	db, err := connect.OpenDB()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	connect.CloseDB(db)
}

func main() {
	CreateTable()
}

package connect

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func CreateMasterKeyTable() *sql.DB {
	createTable := `create table if not exists mastertable (
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		master_key TEXT NOT NULL,
		created_at TEXT,
		updated_at TEXT,
		is_active BOOL
	  );`

	db, err := OpenDB()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	return db
	//CloseDB(db)
}

func InsertMasterKeyDataToDB(db *sql.DB, first_name string, last_name string,
	email string, master_key string, is_active bool) error {
	insertStatement := `INSERT INTO mastertable (first_name, last_name, email, master_key, created_at, updated_at, is_active)
		VALUES($1, $2, $3, $4, $5, $6, $7)`
	time_now := time.Now().Unix()
	created_at := strconv.FormatInt(time_now, 10)
	updated_at := strconv.FormatInt(time_now, 10)
	_, err := db.Exec(insertStatement, first_name, last_name, email, master_key, created_at, updated_at, is_active)
	if err != nil {
		panic(err)
	}
	return err
}

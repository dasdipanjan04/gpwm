package connect

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func CreateMasterKeyTable() *sql.DB {
	createTable := `create table if not exists mastertable (
		id SERIAL PRIMARY KEY,
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
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT email FROM mastertable where(mastertable.email = $3));`
	time_now := time.Now().Unix()
	created_at := strconv.FormatInt(time_now, 10)
	updated_at := strconv.FormatInt(time_now, 10)
	_, err := db.Exec(insertStatement, first_name, last_name, email, master_key, created_at, updated_at, is_active)
	if err != nil {
		panic(err)
	}
	return err
}

func UpdateInfo(db *sql.DB, id int, first_name string, last_name string,
	email string, master_key string, created_at string,
	updated_at string, is_active bool) error {

	updateStatement := `UPDATE mastertable 
	SET first_name = $2, last_name = $3, email = $4, master_key = $5, created_at = $6, updated_at = $7, is_active = $8
	WHERE id = $1;`

	_, err := db.Exec(updateStatement, id, first_name, last_name, email, master_key, created_at, updated_at, is_active)
	if err != nil {
		panic(err)
	}
	return err
}

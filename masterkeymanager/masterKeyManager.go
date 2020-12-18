////////////////////////////////////////////////////////////////////////////////
// Manages various aspects e.g. insert, update, reset of the master key table //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"database/sql"
	"fmt"
	"gpwm/connect"
	"gpwm/masterkeysecure"
	"log"
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
		master_key BYTEA NOT NULL,
		created_at TEXT,
		updated_at TEXT,
		is_active BOOL
	  );`

	db, err := connect.OpenDB()
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return db
}

func InsertMasterKeyDataToDB(db *sql.DB, first_name string, last_name string,
	email string, master_key string, password string, is_active bool) {

	insertStatement := `INSERT INTO mastertable (first_name, last_name, email, master_key, created_at, updated_at, is_active)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT email FROM mastertable where(mastertable.email = $3));`

	time_now := time.Now().Unix()
	created_at := strconv.FormatInt(time_now, 10)
	updated_at := strconv.FormatInt(time_now, 10)

	master_key_byte := []byte(master_key)
	encrypted_master_key := masterkeysecure.EncryptMasterKeyAES(master_key_byte, password)

	_, err := db.Exec(insertStatement, first_name, last_name, email, encrypted_master_key, created_at, updated_at, is_active)
	if err != nil {
		fmt.Println(masterkeysecure.EncryptMasterKeyAES(master_key_byte, password))
		log.Fatalln(err)
		return
	}

	log.Println("Row successfully inserted")
}

func UpdateInfo(db *sql.DB, id int, first_name string, last_name string,
	email string, master_key string, created_at string,
	updated_at string, is_active bool) {

	updateStatement := `UPDATE mastertable 
	SET first_name = $2, last_name = $3, email = $4, master_key = $5, created_at = $6, updated_at = $7, is_active = $8
	WHERE id = $1;`

	_, err := db.Exec(updateStatement, id, first_name, last_name, email, master_key, created_at, updated_at, is_active)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("Successfully Updated.")
}

// Resets master key in the database.
func ResetMasterKey(db *sql.DB, email string, masterKey string) {

	findIdByEmail := fmt.Sprintf(`SELECT id FROM mastertable WHERE email in (%s);`, email)

	id := 0

	err := db.QueryRow(findIdByEmail).Scan(&id)
	if err != nil {
		log.Println(err)
		return
	}

	reserMasterKeyStatement := fmt.Sprintf(`UPDATE mastertable
	SET master_key = (%s)
	WHERE id = (%d);`, masterKey, id)

	_, err = db.Exec(reserMasterKeyStatement)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("You have successfully reset your master key")
}

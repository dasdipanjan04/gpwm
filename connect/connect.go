//////////////////////////////////////////////////////////////
// Creates a successful connection to the postgres database //
//////////////////////////////////////////////////////////////

package connect

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

type PsqlEnv struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func GetPsqlenv() string {

	err := godotenv.Load(path.Join(os.Getenv("HOME"), "go/src/gpwm/connect/psql.env"))
	if err != nil {
		log.Fatalln(err)
	}

	portVal, err := strconv.Atoi(os.Getenv("PSQLPORT"))

	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PSQLHOST"),
		portVal,
		os.Getenv("PSQLUSER"),
		os.Getenv("PSQLPASSWORD"),
		os.Getenv("PSQLDB"))
	return psqlConnect
}

func OpenDB() (*sql.DB, error) {

	psqlEnv := GetPsqlenv()

	// Try open the psql db with the given information.
	psqlDB, connectionError := sql.Open("postgres", psqlEnv)
	if connectionError != nil {
		log.Fatalln(connectionError)
		return nil, connectionError
	}
	return psqlDB, connectionError
}

func ConnectToMasterDB() (*sql.DB, error) {

	psqlDb, err := OpenDB()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	err = psqlDb.Ping()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return psqlDb, err
}

func CloseDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

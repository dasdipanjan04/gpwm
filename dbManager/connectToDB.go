//////////////////////////////////////////////////////////////
// Creates a successful connection to the postgres database //
//////////////////////////////////////////////////////////////

package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PsqlEnv struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func getPsqlenv() PsqlEnv {
	err := godotenv.Load("psql.env")
	if err != nil {
		panic(err)
	}
	var PsqlConnEnv PsqlEnv
	PsqlConnEnv.host = os.Getenv("PSQLHOST")
	portVal, err := strconv.Atoi(os.Getenv("PSQLPORT"))
	PsqlConnEnv.port = portVal
	PsqlConnEnv.user = os.Getenv("PSQLUSER")
	PsqlConnEnv.password = os.Getenv("PSQLPASSWORD")
	PsqlConnEnv.dbname = os.Getenv("PSQLDB")
	return PsqlConnEnv
}

func main() {
	var Envs PsqlEnv = getPsqlenv()

	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Envs.host, Envs.port, Envs.user, Envs.password, Envs.dbname)

	// Try open the psql db with the given information.
	psqlDB, connectionError := sql.Open("postgres", psqlConnect)
	if connectionError != nil {
		panic(connectionError)
	}
	defer psqlDB.Close()

	connectionError = psqlDB.Ping()
	if connectionError != nil {
		panic(connectionError)
	}

	fmt.Println("Postgres connection ok!")
	fmt.Println("Test")
}

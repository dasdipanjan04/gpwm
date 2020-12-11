//////////////////////////////////////////////////////////////
// Creates a successful connection to the postgres database //
//////////////////////////////////////////////////////////////

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dipanjan"
	password = "dipanjan"
	dbname   = "dipanjan"
)

func main() {
	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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

	fmt.Println("Congratulations! Database connected, you're now ready to save passwords")
}

//////////////////////////////////////////////////////////////
// Creates a successful connection to the postgres database //
//////////////////////////////////////////////////////////////

package connect

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/dasdipanjan04/gpwm/helper/glogger"

	"github.com/joho/godotenv"
)

// GetPsqlenv gets the psql env
func GetPsqlenv() string {

	err := godotenv.Load(path.Join(os.Getenv("HOME"), "psql.env"))
	if err != nil {
		glogger.Glog("connect:GetPsqlenv:Load ", err.Error())
		return err.Error()
	}

	portVal, err := strconv.Atoi(os.Getenv("PSQLPORT"))
	if err != nil {
		glogger.Glog("connect:GetPsqlenv:Atoi ", err.Error())
		return err.Error()
	}

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
	psqlDB, err := sql.Open("postgres", psqlEnv)
	if err != nil {
		glogger.Glog("connect:OpenDB:Open ", err.Error())
		return nil, err
	}
	return psqlDB, err
}

func ConnectToMasterDB() (*sql.DB, error) {

	psqlDb, err := OpenDB()
	if err != nil {
		glogger.Glog("connect:ConnectToMasterDB:OpenDB ", err.Error())
		panic(err)
	}

	err = psqlDb.Ping()
	if err != nil {
		glogger.Glog("connect:ConnectToMasterDB:Ping ", err.Error())
		return nil, err
	}

	return psqlDb, err
}

func CloseDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		glogger.Glog("connect:CloseDB:Close ", err.Error())
	}
	return err
}

//////////////////////////////////////////////
// Appends output to a generalized log file //
//////////////////////////////////////////////

package glogger

import (
	"log"
	"os"
)

// Glog logs errors.
func Glog(messageSource string, logMessage string) {

	logFile, err := os.OpenFile("gpwm.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	glogger := log.New(logFile, "gpwm:"+messageSource, log.LstdFlags)
	glogger.Println(logMessage)
}

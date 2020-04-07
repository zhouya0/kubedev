package log

import (
	"fmt"
	"kubedev/pkg/util"
	"log"
	"os"
)

var logPath string = "/var/log/kubedev.log"

func NewLogger() *log.Logger {
	// remove the log file. kubedev is a cli tool, no need to restore history logs.
	if util.CheckExist(logPath) {
		os.Remove(logPath)
	}

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func LogErrorMessage(logger *log.Logger, err error) {
	logger.Println(err.Error())
	fmt.Printf("Error happens, see log file for more details: %s\n", logPath)
}

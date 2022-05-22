package logger

import (
	"log"
	"os"
)

var CommonLog *log.Logger
var ErrorLog *log.Logger

func init() {
	logFile, err := os.OpenFile("../../internal/logger/log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// 0666 is the code to set the default file permission to  modify the file
	if err != nil {
		log.Println("Error opening file ", err)
		os.Exit(1)
	}

	CommonLog = log.New(logFile, "Common Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(logFile, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

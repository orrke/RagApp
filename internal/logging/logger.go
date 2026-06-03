package logging

import (
	"log"
	"os"
	"path"
	"time"
)

var Logger *log.Logger
var LoggerPath string

// LogSetup is the function responsible for setting up the logging system of the server.
func LogSetup() (*os.File, error) {
	filePath := getLogFileName()

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	return file, nil //return the file so that it can be closed at the end of main.
}

// getLogFileName is the function to get the filename (depends on the time)
func getLogFileName() string {
	currentTime := time.Now()

	return path.Join(LoggerPath, currentTime.Format("20060102T150405Z")+".log")
}

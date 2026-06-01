package logging

import (
	"log"
	"os"
	"path"
	"time"
)

var Logger *log.Logger
var LoggerPath string

func LogSetup() error {
	filePath := getLogFileName()

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	return nil
}

func getLogFileName() string {
	currentTime := time.Now()

	path.Join(LoggerPath, currentTime.Format(time.RFC3339)+".log")
	return "" + currentTime.Format(time.UnixDate) + ".log"
}

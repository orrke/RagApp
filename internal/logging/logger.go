package logging

import (
	"log"
	"os"
	"path"
	"time"
)

var Logger *log.Logger
var LoggerPath string

func LogSetup() (*os.File, error) {
	filePath := getLogFileName()

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	return file, nil
}

func getLogFileName() string {
	currentTime := time.Now()

	return path.Join(LoggerPath, currentTime.Format("20060102T150405Z")+".log")
}

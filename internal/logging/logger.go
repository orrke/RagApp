package logging

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"
)

var Logger *log.Logger
var LoggerPath string
var MaxLogCount = 10

// LogSetup is the function responsible for setting up the logging system of the server.
func LogSetup() (*os.File, error) {
	filePath := getLogFileName()

	err := rollLogFiles(LoggerPath)
	if err != nil {
		return nil, err
	}

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

func rollLogFiles(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if len(entries) < 10 {
		return nil
	}

	slices.SortFunc(entries, func(a, b os.DirEntry) int {
		aInfo, _ := a.Info()
		bInfo, _ := b.Info()
		return aInfo.ModTime().Compare(bInfo.ModTime())
	})

	for len(entries) > MaxLogCount {
		err = os.Remove(filepath.Join(LoggerPath, entries[0].Name()))
		if err != nil {
			return err
		}
		entries = entries[1:]
	}

	return nil
}

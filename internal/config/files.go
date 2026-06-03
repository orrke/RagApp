package config

import (
	"os"
	"path/filepath"
)

// SetupFolders is the function responsible for the creation of the folders the app will need to use for its config and log files.
// It can't use logging since it's set up before the logger itself.
func SetupFolders(root string) error {
	//Create the root folder
	err := os.MkdirAll(root, os.ModePerm)
	if err != nil {
		return err
	}

	//Create the logs folder
	err = os.MkdirAll(filepath.Join(root, "logs"), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

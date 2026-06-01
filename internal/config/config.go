package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

// GetConfigFromFile replaces the Config variable with the content of the configuration file at Path
func GetConfigFromFile() error {
	//Lock the config down
	Lock.Lock()

	configPath := path.Join(Path, "server_config.json")

	if fileExists(configPath) {
		//Open the configuration file, it's plain text so we don't need an external crate
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println("aaaaa")
			return err
		}
		defer file.Close()

		//decode the config file into the ServerConfig struct
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&Config)

		Lock.Unlock()
		return nil
	}

	//file doesn't exist, get the default config and save it to the path
	defaultConfig := ServerConfig{}.Default()
	Config = defaultConfig

	lang := Config.Language
	if _, exists := Languages[lang]; !exists {
		return errors.New("Unsupported language: " + lang)
	}

	Lock.Unlock()
	return SaveConfigToFile()
}

// SaveConfigToFile saves the current content of the ServerConfig struct into the config file at Path
func SaveConfigToFile() error {
	//Locks the Config struct down
	Lock.RLock()
	defer Lock.RUnlock()

	configPath := path.Join(Path, "server_config.json")

	//Open the file with write and create permissions
	//0644: owner can rw, groups can r and others can r
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	//Marshal the struct into a json
	b, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		return err
	}

	//Write the JSON into the file as plain text
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// StoreCurrentDate save the current time to the last_update field of the json config file
func StoreCurrentDate() error {
	Lock.Lock()
	Config.LastUpdate = new(time.Now())
	Lock.Unlock()

	return SaveConfigToFile()
}

// fileExists checks if a file is present at the specified path
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

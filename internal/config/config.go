package config

import (
	"RagApp/internal/logging"
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"
)

// GetConfigFromFile replaces the Config variable with the content of the configuration file at Path
func GetConfigFromFile() error {
	logging.Trace("GetConfigFromFile")

	//Lock the config down
	logging.Debug("Locking down config")
	Lock.Lock()

	configPath := path.Join(Path, "server_config.json")

	if fileExists(configPath) {
		//Open the configuration file, it's plain text so we don't need an external crate
		file, err := os.Open(configPath)
		if err != nil {
			logging.Debug("Config file not found")
			logging.Trace("Returning GetConfigFromFile")
			return err
		}
		defer file.Close()

		//decode the config file into the ServerConfig struct
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&Config)

		Lock.Unlock()
		logging.Debug("Config released")
		logging.Debug("File found")
		logging.Trace("returning GetConfigFromFile")
		return nil
	}

	//file doesn't exist, get the default config and save it to the path
	defaultConfig := ServerConfig{}.Default()
	Config = defaultConfig

	lang := Config.Language
	if _, exists := Languages[lang]; !exists {
		logging.Debug("Unsupported language found!")
		return errors.New("Unsupported language: " + lang)
	}

	Lock.Unlock()
	logging.Debug("Config lock released")
	logging.Debug("File not found")
	logging.Trace("returning GetConfigFromFile")
	return SaveConfigToFile()
}

// SaveConfigToFile saves the current content of the ServerConfig struct into the config file at Path
func SaveConfigToFile() error {
	logging.Trace("SaveConfigToFile")

	//Locks the Config struct down
	logging.Debug("Locking down config")
	Lock.RLock()
	defer func() {
		Lock.RUnlock()
		logging.Debug("Releasing config")
	}()

	configPath := path.Join(Path, "server_config.json")

	//Open the file with write and create permissions
	//0644: owner can rw, groups can r and others can r
	logging.Debug("Saving config")
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logging.Trace("returning SaveConfigToFile")
		return err
	}
	defer file.Close()

	//Marshal the struct into a json
	b, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		logging.Trace("returning SaveConfigToFile")
		return err
	}

	//Write the JSON into the file as plain text
	logging.Debug("Writing config to file")
	_, err = file.Write(b)
	if err != nil {
		logging.Trace("returning SaveConfigToFile")
		return err
	}

	logging.Trace("returning SaveConfigToFile")
	return nil
}

// StoreCurrentDate save the current time to the last_update field of the json config file
func StoreCurrentDate() error {
	logging.Trace("StoreCurrentDate")

	logging.Debug("Locking down config")
	Lock.Lock()
	Config.LastUpdate = new(time.Now())
	Lock.Unlock()
	logging.Debug("Config released")

	err := SaveConfigToFile()

	if err != nil {
		logging.Error("Unable to save config file: " + err.Error())
		logging.Trace("returning StoreCurrentDate")
		return err
	}

	logging.Trace("returning StoreCurrentDate")
	return nil
}

// fileExists checks if a file is present at the specified path
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

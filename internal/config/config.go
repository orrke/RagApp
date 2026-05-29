package config

import (
	"encoding/json"
	"os"
	"time"
)

const (
	Path = "C:\\RAG\\config\\server_config.json"
)

// GetConfigFromFile replaces the Config variable with the content of the configuration file at Path
func GetConfigFromFile() error {
	//Lock the config down
	Lock.Lock()
	defer Lock.Unlock() //Unlocks when the function is done

	//Open the configuration file, it's plain text so we don't need an external crate
	file, err := os.Open(Path)
	if err != nil {
		return err
	}
	defer file.Close()

	//decode the config file into the ServerConfig struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)

	return nil
}

// SaveConfigToFile saves the current content of the ServerConfig struct into the config file at Path
func SaveConfigToFile() error {
	//Locks the Config struct down
	Lock.RLock()
	defer Lock.RUnlock()

	//Open the file with write permissions
	file, err := os.OpenFile(Path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	//Marshal the struct into a json
	b, err := json.Marshal(Config)
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

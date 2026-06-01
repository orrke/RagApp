package config

import (
	"RagApp/internal/logging"
	"time"
)

// ServerConfig is the struct that holds all the parameter the server needs to function
// It's the struct equivalent of the JSON file at Path
type ServerConfig struct {
	BleveIndexPath string     `json:"bleve_index_path"`
	DocsPath       string     `json:"docs_path"`
	Model          string     `json:"model"`
	Language       string     `json:"language"`
	LastUpdate     *time.Time `json:"last_update"`
}

func (s ServerConfig) Default() ServerConfig {
	logging.Trace("Defaulting server config")

	return ServerConfig{
		BleveIndexPath: Path, //gets index.bleve added to it, which contains the actual data
		DocsPath:       "",   //TODO: put that into a variable at server start
		Model:          "gemma4:latest",
		Language:       "en",
		LastUpdate:     nil,
	}
}

// SetArgs sets the given cli arguments to the given config object, and saves the main Config to a file (assuming it's the main config that's getting updated)
func (s ServerConfig) SetArgs(docsPath *string, model *string) error {
	logging.Trace("Setting arguments")

	logging.Debug("Locking down config")
	Lock.Lock()

	if docsPath != nil {
		s.DocsPath = *docsPath
	}

	if model != nil {
		s.Model = *model
	}

	Lock.Unlock()
	logging.Debug("Releasing config lock")

	err := SaveConfigToFile()
	if err != nil {
		return err
	}

	logging.Trace("returning SaveConfigToFile")
	return nil
}

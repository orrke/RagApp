package config

import "time"

// ServerConfig is the struct that holds all the parameter the server needs to function
// It's the struct equivalent of the JSON file at Path
type ServerConfig struct {
	BleveIndexPath string     `json:"bleve_index_path"`
	DocsPath       string     `json:"docs_path"`
	Model          string     `json:"model"`
	Language       string     `json:"language"`
	LastUpdate     *time.Time `json:"last_update"`
}

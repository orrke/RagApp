package config

import (
	"sync"
)

// Global configuration variables
var (
	Path   string
	Lock   sync.RWMutex
	Config ServerConfig
)

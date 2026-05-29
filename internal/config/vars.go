package config

import (
	"sync"
)

// Global configuration variables
var (
	Lock   sync.RWMutex
	Config ServerConfig
)

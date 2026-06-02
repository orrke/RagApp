package main

import (
	"RagApp/internal/api"
	"RagApp/internal/config"
	"RagApp/internal/database"
	"RagApp/internal/logging"
	"RagApp/web"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	//Get the config directory for the OS
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	//Put the configuration file in the default config directory
	config.Path = path.Join(configDir, "RagApp")
	logging.LoggerPath = path.Join(config.Path, "logs")

	//Set up the logger
	file, err := logging.LogSetup()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Get the command line arguments for the server (address and port)
	serverArgs := GetArgs()

	//Load the configuration file into config.Config (global variable), access with the config.Lock RWLock
	err = config.GetConfigFromFile()
	if err != nil {
		logging.Fatal(err.Error())
	}

	//Get some variables we'll use to init the server, from the config file we just loaded
	logging.Debug("Locking down config")
	config.Lock.RLock()
	docsPath := config.Config.DocsPath
	time := config.Config.LastUpdate
	config.Lock.RUnlock()
	logging.Debug("Releasing config lock")

	//Set the variables if they're not nil
	err = config.Config.SetArgs(serverArgs.docsPath, serverArgs.model)
	if err != nil {
		log.Fatal(err)
	}

	//docsPath isn't set by default, meaning the config file is the default one
	if docsPath == "" {
		log.Fatal("docs path is required if no config file is provided. Provide it with the --docs flag")
	}

	//Set the index to the default configuration file value
	err = database.SetIndex()
	if err != nil {
		log.Fatal(err)
	}

	//Store all the files in the directory specified in the config file inside the bleve index
	err = database.StoreAllFilesInDefaultDir(docsPath, time)
	if err != nil {
		log.Fatal(err)
	}

	//Setup server handler routes
	web.SetupRoutes()
	api.SetupRoutes("/api/v1")

	//Start server
	err = http.ListenAndServe(serverArgs.address+":"+serverArgs.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

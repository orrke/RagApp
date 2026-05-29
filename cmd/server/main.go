package main

import (
	"RagApp/internal/api"
	"RagApp/internal/config"
	"RagApp/internal/database"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	//Gte the config directory for the OS
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	//Put the configuration file in the default config directory
	config.Path = path.Join(configDir, "RagApp")

	//Get the command line arguments for the server (address and port)
	serverArgs := GetArgs()

	//Load the configuration file into config.Config (global variable), access with the config.Lock RWLock
	err = config.GetConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}

	//Get some variables we'll use to init the server, from the config file we just loaded
	config.Lock.RLock()
	docsPath := config.Config.DocsPath
	time := config.Config.LastUpdate
	config.Lock.RUnlock()

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
	http.HandleFunc("/", api.ServerIsAlive)        //GET
	http.HandleFunc("/search", api.SearchDocument) //POST
	http.HandleFunc("/reset", api.ResetIndex)      //GET
	http.HandleFunc("/reconfig", api.ReloadConfig) //GET

	//Start server
	err = http.ListenAndServe(serverArgs.address+":"+serverArgs.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

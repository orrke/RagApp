package main

import (
	"RagApp/internal/api"
	"RagApp/internal/config"
	"RagApp/internal/database"
	"log"
	"net/http"
)

func main() {
	//Get the command line arguments for the server (address and port)
	serverArgs := GetArgs()

	//Load the configuration file into config.Config (global variable), access with the config.Lock RWLock
	err := config.GetConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}

	//Get some variables we'll use to init the server, from the config file we just loaded
	config.Lock.RLock()
	docsPath := config.Config.DocsPath
	time := config.Config.LastUpdate
	config.Lock.RUnlock()

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

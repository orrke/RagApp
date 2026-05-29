package database

import (
	"RagApp/internal/config"
	"os"
	"path"
	"sync"

	"github.com/blevesearch/bleve/v2"
)

var (
	Index     bleve.Index
	IndexLock sync.RWMutex
)

// SetIndex sets the Index global variable to the index specified in the config file
func SetIndex() error {
	config.Lock.RLock()
	bleveIndexPath := config.Config.BleveIndexPath
	config.Lock.RUnlock()

	//get the index from the default path
	index, err := GetOrCreateIndex(bleveIndexPath)
	if err != nil {
		return err
	}

	//Set the index global variable
	IndexLock.Lock()
	defer IndexLock.Unlock()
	Index = index

	return nil
}

func ResetIndex() error {
	IndexLock.Lock()
	defer IndexLock.Unlock()

	//Fetch some useful config variables
	config.Lock.RLock()
	bleveIndexPath := config.Config.BleveIndexPath
	docsPath := config.Config.DocsPath
	time := config.Config.LastUpdate
	config.Lock.RUnlock()

	//Flush the entire bleve index
	fullPath := path.Join(bleveIndexPath, "index.bleve")
	err := os.RemoveAll(fullPath)
	if err != nil {
		return err
	}

	//Recreate the index
	//(it's the wrapper but in reality there's no index existing so it'll create a new one)
	index, err := GetOrCreateIndex(bleveIndexPath)
	if err != nil {
		return err
	}

	Index = index

	//Store all the files in the default directory inside the index
	err = StoreAllFilesInDefaultDir(docsPath, time)
	if err != nil {
		return err
	}

	return nil
}

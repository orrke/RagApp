package database

import (
	"RagApp/internal/config"
	"RagApp/internal/logging"
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
	logging.Trace("SetIndex")
	logging.Debug("Locking down config")
	config.Lock.RLock()
	bleveIndexPath := config.Config.BleveIndexPath
	config.Lock.RUnlock()
	logging.Debug("Releasing index lock")

	//get the index from the default path
	index, err := GetOrCreateIndex(bleveIndexPath)
	if err != nil {
		return err
	}

	//Set the index global variable
	logging.Debug("Locking down index")
	IndexLock.Lock()
	defer func() {
		IndexLock.Unlock()
		logging.Debug("Releasing index lock")
	}()
	Index = index

	logging.Trace("returning SetIndex")
	return nil
}

func ResetIndex() error {
	logging.Trace("ResetIndex")
	logging.Debug("Locking down config")
	IndexLock.Lock()
	defer func() {
		IndexLock.Unlock()
		logging.Debug("Releasing index lock")
	}()

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

	logging.Trace("returning ResetIndex")
	return nil
}

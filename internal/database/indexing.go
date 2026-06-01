package database

import (
	"RagApp/internal/config"
	"RagApp/internal/logging"
	"encoding/json"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/fr"
)

// Document :
//
// The document we get from the file we extract directly
type Document struct {
	Title   string
	Content string
	Path    string
}

// StoredDocument :
//
// The document we give to bleve
// there's no title since it's in the path
// the content of the file is in the JSON format
// so that it can easily be stored and retrieved from the database
// so we can give it to the AI after that
type StoredDocument struct {
	Content string
	Path    string
	Raw     string
}

// Type is a function required by bleve for indexing the struct
func (d Document) Type() string {
	return "myDoc"
}

// createIndex creates the bleve index at the specified path.
// The index creates a mapping for the content in the specified language
// It also creates a mapping for the path, searching by keyword
// It creates a mapping for the raw field, but it doesn't use it to searchfor keywords
func createIndex(bleveIndexPath string) (bleve.Index, error) {
	logging.Trace("createIndex")

	logging.Debug("Locking down config")
	config.Lock.RLock()
	language := config.Config.Language
	config.Lock.RUnlock()
	logging.Debug("Releasing config lock")

	//create the index mapping and the documents mapping
	indexMapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()

	//Create the mapping for the contents of the file
	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Analyzer = language

	//Create the mapping for the file path
	pathMapping := bleve.NewTextFieldMapping()
	pathMapping.Analyzer = "keyword"

	//Create the mapping for the raw field
	rawMapping := bleve.NewTextFieldMapping()
	rawMapping.Store = true  //store the file in the index
	rawMapping.Index = false //don't use the mapping to search for given keywords

	//add the field mappings to the document mapping
	docMapping.AddFieldMappingsAt("Path", pathMapping)
	docMapping.AddFieldMappingsAt("Content", contentMapping)
	docMapping.AddFieldMappingsAt("Raw", rawMapping)

	//add the mapping to the index
	indexMapping.AddDocumentMapping("myDoc", docMapping)

	//return the index
	bleveIndex, err := bleve.New(bleveIndexPath+"index.bleve", indexMapping)
	if err != nil {
		return nil, err
	}

	logging.Trace("returning createIndex")
	return bleveIndex, nil
}

// getIndex fetches an existing index at the specified path
func getIndex(bleveIndexPath string) (bleve.Index, error) {
	logging.Trace("getIndex")

	index, err := bleve.Open(bleveIndexPath + "index.bleve")
	if err != nil {
		return nil, err
	}

	logging.Trace("returning getIndex")
	return index, nil
}

// GetOrCreateIndex is a unified wrapper around the createIndex and getIndex functions
func GetOrCreateIndex(bleveIndexPath string) (bleve.Index, error) {
	logging.Trace("GetOrCreateIndex")

	index, err := getIndex(bleveIndexPath)
	if err != nil {
		index, err = createIndex(bleveIndexPath)
		if err != nil {
			return nil, err
		}
	}

	logging.Trace("returning GetOrCreateIndex")
	return index, nil
}

// addDocumentToIndex adds a single document to the given index
func addDocumentToIndex(index bleve.Index, document Document) error {
	logging.Trace("addDocumentToIndex")

	//marshal the original doc into the raw json format
	rawDoc, err := json.Marshal(document)
	if err != nil {
		return err
	}
	//create the document we'll actually store inside the index
	doc := StoredDocument{Content: document.Content, Path: document.Path, Raw: string(rawDoc)}

	//store the document into the index
	//the primary key is the path, since we only store a single file once
	logging.Debug("locking down index")
	IndexLock.Lock()
	err = index.Index(document.Path, doc)
	IndexLock.Unlock()
	logging.Debug("Releasing index lock")
	if err != nil {
		return err
	}

	logging.Trace("returning addDocumentToIndex")
	return nil
}

// AddAllDocumentsToIndex adds all the documents in the given array into the specified index
func AddAllDocumentsToIndex(index bleve.Index, documents []Document) error {
	logging.Trace("AddAllDocumentsToIndex")
	for _, document := range documents {
		err := addDocumentToIndex(index, document)
		if err != nil {
			return err
		}
	}
	logging.Trace("returning AddAllDocumentsToIndex")
	return nil
}

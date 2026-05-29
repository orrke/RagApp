package database

import (
	"RagApp/internal/config"
	"RagApp/internal/extract"
	"os"
	"path/filepath"
	"time"
)

// GetAllFilesInDir :
//
// Retrieve every single file in the specified directory, depending on the last database update time
func GetAllFilesInDir(dir string, lastUpdate *time.Time) ([]Document, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []Document

	//iterate through directory children
	for _, entry := range entries {
		if entry.IsDir() {
			//child is a directory, we iterate through it too
			documents, err := GetAllFilesInDir(filepath.Join(dir, entry.Name()), lastUpdate)
			if err != nil {
				return nil, err
			}
			result = append(result, documents...)
		} else {
			//child is a file, we need ot get the modification date before doing anything
			info, err := entry.Info()
			if err != nil {
				return nil, err
			}

			if lastUpdate != nil && info.ModTime().Before(*lastUpdate) {
				//File is already indexed, skipping
				continue
			}

			//Get the content of the document as text from the extract package
			content, err := extract.Extract(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}

			//Append the fetched document into the result slice
			result = append(result, Document{
				Path:    filepath.Join(dir, entry.Name()),
				Title:   entry.Name(),
				Content: content,
			})
		}
	}

	return result, nil
}

// StoreAllFilesInDefaultDir :
//
// Stores all the files in the default directory inside the bleve index.
func StoreAllFilesInDefaultDir(defaultDocsPath string, lastUpdate *time.Time) error {
	//fetches all the files in the specified directory
	files, err := GetAllFilesInDir(defaultDocsPath, lastUpdate)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		//No files to be updated, just return, don't even need to update the last_update time
		return nil
	}

	//Update the update time to right now
	err = config.StoreCurrentDate()
	if err != nil {
		return err
	}

	//Add all the fetched documents to the index
	return AddAllDocumentsToIndex(Index, files)
}

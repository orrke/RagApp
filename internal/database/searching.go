package database

import (
	"encoding/json"

	"github.com/blevesearch/bleve/v2"
)

// SearchDatabase searches for documents with the given query through the specified index
// It returns documents sorted by their score in descending order
func SearchDatabase(query string, index bleve.Index) ([]Document, error) {
	bleveQuery := bleve.NewMatchQuery(query)
	searchRequest := bleve.NewSearchRequest(bleveQuery)
	searchRequest.SortBy([]string{"-_score"}) //Sort by score descending
	searchRequest.Fields = []string{"Raw"}    //Only fetch the 'Raw' field

	searchResults, err := index.Search(searchRequest) //search through the database
	if err != nil {
		return nil, err
	}

	var results []Document

	for _, hit := range searchResults.Hits {
		//Get the 'Raw' field from the result
		rawStr := hit.Fields["Raw"].(string)

		//Recreate the document from the Raw field, which is in a JSON format
		var document Document
		err := json.Unmarshal([]byte(rawStr), &document)
		if err != nil {
			return nil, err
		}

		results = append(results, document)
	}

	return results, nil
}

// Document.String converts a document into a string
func (d Document) String() string {
	return "\n-- Start of Document -- Title: " + d.Title + " Path: " + d.Path + " Content: " + d.Content + " -- End of Document --\n"
}

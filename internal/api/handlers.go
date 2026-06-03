package api

import (
	"RagApp/internal/config"
	"RagApp/internal/database"
	"RagApp/internal/logging"
	"RagApp/internal/ollama"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ServerIsAlive is the root (/) handler, returns OK to confirm that the server is indeed alive.
func ServerIsAlive(w http.ResponseWriter, _ *http.Request) {
	logging.Info("Received IsAlive request")
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SearchDocument is the handler for the /search endpoint, the main endpoint of the server.
func SearchDocument(w http.ResponseWriter, r *http.Request) {
	logging.Info("Received SearchDocument request")

	if r.Method != http.MethodPost { //Only allow POST
		logging.Info("Invalid request method: " + r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Unmarshal the data into the SearchParameters struct
	//(which is actually a SearchParameter right now but if I decide to add parameters it'll be accurate)
	var data SearchParameters
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logging.Info("Error decoding search parameters: " + err.Error())
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	//Get some useful variables from the config file
	logging.Debug("Locking down config")
	config.Lock.RLock()
	model := config.Config.Model
	language := config.Config.Language
	config.Lock.RUnlock()
	logging.Debug("Releasing config Lock")

	//Query the model for the necessary documents
	res, err := ollama.AskModel(data.Query, database.Index, model, language)
	if err != nil {
		logging.Info("Error asking model: " + err.Error())

		//Formulate the failure response
		response := SearchResponse{
			Result:  "failure",
			Content: "Failed to ask model: " + err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	logging.Trace("Successfully asked model")

	//Formulate the success response
	response := SearchResponse{
		Result:  "success",
		Content: res,
	}

	//Make the header to return to the user
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logging.Info("Error encoding search response: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SearchRawDocument is the handler for getting documents without any model processing.
func SearchRawDocument(w http.ResponseWriter, r *http.Request) {
	logging.Info("Received SearchRawDocument request")

	if r.Method != http.MethodPost {
		logging.Info("Invalid request method: " + r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data SearchParameters
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logging.Info("Error decoding search parameters: " + err.Error())
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	res, err := database.SearchDatabase(data.Query, database.Index)
	if err != nil {
		logging.Info("Error searching database: " + err.Error())

		response := SearchResponse{
			Result:  "failure",
			Content: "Failed to search database: " + err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if len(res) == 0 {
		response := SearchResponse{
			Result:  "failure",
			Content: "No results found",
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var finalRes []string
	for _, doc := range res {
		finalRes = append(finalRes, doc.String())
	}

	w.Header().Set("Content-Type", "application/json")
	response := SearchResponse{
		Result:  "success",
		Content: strings.Join(finalRes, ","),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateIndex is the handler responsible for updating the index with no reset happening.
func UpdateIndex(w http.ResponseWriter, r *http.Request) {
	logging.Info("Received UpdateIndex request")

	config.Lock.RLock()
	docsPath := config.Config.DocsPath
	lastUpdate := config.Config.LastUpdate
	config.Lock.RUnlock()

	err := database.StoreAllFilesInDefaultDir(docsPath, lastUpdate)
	if err != nil {
		logging.Info("Error updating index: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ResetIndex is the handler for the /reset endpoint
func ResetIndex(w http.ResponseWriter, _ *http.Request) {
	logging.Info("Received ResetIndex request")

	//Reset the entire index
	err := database.ResetIndex()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ReloadConfig is the handler for the /reconfig endpoint
func ReloadConfig(w http.ResponseWriter, _ *http.Request) {
	logging.Info("Received ReloadConfig request")

	//Reload the configuration from the json file
	err := config.GetConfigFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

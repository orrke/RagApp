package api

import (
	"RagApp/internal/config"
	"RagApp/internal/database"
	"RagApp/internal/ollama"
	"encoding/json"
	"fmt"
	"net/http"
)

// ServerIsAlive is the root (/) handler, returns OK to confirm that the server is indeed alive.
func ServerIsAlive(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SearchDocument is the handler for the /search endpoint, the main endpoint of the server.
func SearchDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //Only allow POST
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Unmarshal the data into the SearchParameters struct
	//(which is actually a SearchParameter right now but if I decide to add parameters it'll be accurate)
	var data SearchParameters
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	//Get some useful variables from the config file
	config.Lock.RLock()
	model := config.Config.Model
	language := config.Config.Language
	config.Lock.RUnlock()

	//Query the model for the necessary documents
	res, err := ollama.AskModel(data.Query, database.Index, model, language)
	if err != nil {
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
	}

	//Formulate the success response
	response := SearchResponse{
		Result:  "success",
		Content: res,
	}

	//Make the header to return to the user
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ResetIndex is the handler for the /reset endpoint
func ResetIndex(w http.ResponseWriter, _ *http.Request) {
	//Reset the entire index
	err := database.ResetIndex()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ReloadConfig is the handler for the /reconfig endpoint
func ReloadConfig(w http.ResponseWriter, _ *http.Request) {
	//Reload the configuration from the json file
	err := config.GetConfigFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

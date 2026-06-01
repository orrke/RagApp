package ollama

import (
	"RagApp/internal/logging"
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	ollamaUrl = "http://localhost:11434/api/chat"
)

// queryModel does the actual request to the ollama endpoint
func queryModel(request NewChatRequest) (*ChatResponse, error) {
	logging.Trace("queryModel")

	//Marshal the request struct into a json
	jsonData, err := json.Marshal(request.ToChatRequest())
	if err != nil {
		return nil, err
	}

	//Do the post request to ollama
	logging.Info("Doing POST request to ollama server")
	post, err := http.Post(ollamaUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()

	//infinite loop to ignore the thinking part
	for {
		var chunk ChatResponse

		//decode the result into a ChatResponse struct
		if err := json.NewDecoder(post.Body).Decode(&chunk); err != nil {
			return nil, err
		}
		logging.Info("Receiving chunk from ollama server")

		//Model is done thinking, we can process the result
		if chunk.Done {
			logging.Trace("returnng queryModel")
			return &chunk, nil
		}
	}
}

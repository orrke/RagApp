package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	ollamaUrl = "http://localhost:11434/api/chat"
)

// queryModel does the actual request to the ollama endpoint
func queryModel(request NewChatRequest) (*ChatResponse, error) {
	//Marshal the request struct into a json
	jsonData, err := json.Marshal(request.ToChatRequest())
	if err != nil {
		return nil, err
	}

	//Do the post request to ollama
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

		//Model is done thinking, we can process the result
		if chunk.Done {
			return &chunk, nil
		}
	}
}

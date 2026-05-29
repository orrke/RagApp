package ollama

import (
	"github.com/blevesearch/bleve/v2"
)

// AskModel is a function to query a model through ollama to get information from the bleve index
func AskModel(query string, index bleve.Index, model string, language string) (string, error) {
	//variables for iteration and result
	queryingModel := true
	var modelAnswer string

	//Get the initial request, made up of the 2 initial messages:
	//The system message, containing basic model instructions
	//The user message, containing the user's query
	currentRequest, messageHistory := GetInitialRequest(query, model, language)

	for queryingModel {
		//Query the model and get the response
		response, err := queryModel(currentRequest)
		if err != nil {
			return "", err
		}

		if len(response.Message.ToolCalls) > 0 {
			//add the message to the history,
			//since ollama doesn't keep an internal state of the conversation
			messageHistory = append(messageHistory, response.Message)

			//Handle the model's tool call
			requestPointer, err := handleToolCall(*response, &messageHistory, index)

			if err != nil {
				return "", err
			}
			currentRequest = *requestPointer
		} else {
			//No tool is being called, which means we probably have an answer from the model
			modelAnswer = response.Message.Content
			queryingModel = false
		}
	}

	return modelAnswer, nil
}

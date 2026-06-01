package ollama

import "RagApp/internal/config"

// CreateSystemPrompt creates the initial system prompt for the model
func CreateSystemPrompt(language string) string {
	return "Answer the user in the following language only: " + config.Languages[language] + "\n" +
		"You are an AI agent used for Research Advanced Generation.\n" +
		"You have access to tools to help you answer the user.\n" +
		"Do not use any context coming from any other source than the main tool you have access to.\n" +
		"You can use the \"query_database\" tool in order to search for things in the database.\n" +
		"The user will ask a question, from that you can call the query_database tool to answer them.\n" +
		"The tool searches in a Bleve database in go, formulate your request for that database.\n" +
		"The files can be quite big, as such you should not take too many files at once (I would recommend not more than 5)\n" +
		"Give a very brief answer containing only the content the user is searching for.\n" +
		"If you have access to a source, give its full path alongside the answer you gave to the user."
}

// CreateUserPrompt is a function to make the initial user prompt
func CreateUserPrompt(query string, language string) string {
	return "Use the tools at your disposition to answer the user in the following language only: " + config.Languages[language] + "\n" +
		"User query: " + query + "\n"
}

// GetInitialRequest is a function to create the initial NewChatRequest object, as well as the beginning of the message history
func GetInitialRequest(query string, model string, language string) (NewChatRequest, []Message) {
	//get the system and user prompt
	systemPrompt := CreateSystemPrompt(language)
	userPrompt := CreateUserPrompt(query, language)

	//Create the first messages with the right roles
	systemMessage := Message{
		Role:      "system",
		Content:   systemPrompt,
		ToolCalls: nil,
	}
	userMessage := Message{
		Role:      "user",
		Content:   userPrompt,
		ToolCalls: nil,
	}
	messageHistory := []Message{systemMessage, userMessage}

	//Create the request
	return NewChatRequest{
		Model:    model,
		Messages: messageHistory,
		Tools:    []Tool{getQueryTool()},
	}, messageHistory
}

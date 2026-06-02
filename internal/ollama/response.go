package ollama

import (
	"RagApp/internal/database"
	"RagApp/internal/logging"
	"errors"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

// handleToolCall is a function to handle a model's tool call and call the correct function
func handleToolCall(response ChatResponse, history *[]Message, index bleve.Index) (*NewChatRequest, error) {
	logging.Trace("handleToolCall")

	if len(response.Message.ToolCalls) < 1 {
		return nil, errors.New("no tool calls found")
	}

	//iterate through the tool calls
	for _, call := range response.Message.ToolCalls {
		switch call.Function.Name {
		case "query_database":
			{
				//called the query tool
				//If I had a bit more tools I'd put them in their own functions
				//But in this case it's the only tool so it's really not that bad

				var start int
				var end int

				//get the arguments
				query := call.Function.Arguments["query"].(string)
				startIndex := call.Function.Arguments["start"]
				endIndex := call.Function.Arguments["end"]

				//search the database for the query
				results, err := database.SearchDatabase(query, index)
				if err != nil {
					return nil, err
				}

				if len(results) < 1 {
					message := Message{
						Role:    "tool",
						Content: "No match found in Bleve index.",
					}

					*history = append(*history, message)
					continue
				}

				//set the default values for endIndex and startIndex
				switch v := startIndex.(type) {
				case nil:
					start = 0
				case string:
					start, _ = strconv.Atoi(v)
				default:
					start = int(v.(float64))
				}

				start = max(start, 0)

				switch v := endIndex.(type) {
				case nil:
					end = 2
				case string:
					end, _ = strconv.Atoi(v)
				default:
					end = int(v.(float64))
				}

				var res []string
				//clip the value of the end index in between the start value and the amount of results (minus one since indices start at 0)
				stopAt := max(min(len(results)-1, end), start)
				for i := start; i <= stopAt; i++ {
					res = append(res, results[i].String())
				}

				//make the tool message
				message := Message{
					Role:    "tool",
					Content: strings.Join(res, "\n"),
				}

				//add the message to the history (they're pointers otherwise I can't modify the history)
				*history = append(*history, message)
			}
		default:
			return nil, errors.New("Unsupported tool called: " + call.Function.Name) //model called another tool for some reason
		}
	}

	logging.Trace("returning handleToolCall")
	//formulate the request
	return &NewChatRequest{
		Model:    response.Model,
		Messages: *history,
		Tools:    []Tool{getQueryTool()},
	}, nil
}

// getQueryTool is a function to generate the query tool for the model
func getQueryTool() Tool {
	return Tool{
		Type: "function", //always the same I think
		Function: FunctionDef{
			Name:        "query_database",
			Description: "Query database for tool calls",
			Parameters: map[string]interface{}{
				"type":     "object",
				"required": []string{"query"},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "The query to give to bleve for the search",
				},
				"start": map[string]interface{}{
					"type":        "number",
					"description": "The first document to search for starting from 0. defaults to 0",
				},
				"end": map[string]interface{}{
					"type":        "number",
					"description": "The last document to search for starting from 0. defaults to 2 or the amount of documents retrieved",
				},
			},
		},
	}
}

package ollama

// ChatRequest represents the request body for the /api/chat endpoint
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role      string     `json:"role"`    // "user", "assistant", "system" or "tool"
	Content   string     `json:"content"` // The message content
	Thinking  string     `json:"thinking,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a function call made by the model
type ToolCall struct {
	ID       string       `json:"id"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents the function details in a tool call
type FunctionCall struct {
	Index     int                    `json:"index"`
	Name      string                 `json:"name"` // Name of the function to call
	Arguments map[string]interface{} `json:"arguments"`
}

// Tool represents a function definition for function calling
type Tool struct {
	Type     string      `json:"type"`     // "function"
	Function FunctionDef `json:"function"` // Function definition
}

// FunctionDef represents the function definition
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"` // JSON Schema
}

// ChatResponse represents the response from the /api/chat endpoint
type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

// NewChatRequest creates a new ChatRequest with streaming disabled
type NewChatRequest struct {
	Model    string
	Messages []Message
	Tools    []Tool
}

// ToChatRequest converts to ChatRequest with stream disabled
func (r *NewChatRequest) ToChatRequest() *ChatRequest {
	return &ChatRequest{
		Model:    r.Model,
		Messages: r.Messages,
		Stream:   false,
		Tools:    r.Tools,
	}
}

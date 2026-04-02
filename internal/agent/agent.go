package agent

import (
	"context"
	"pedals/internal/config"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the chat completion API
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a response from the chat completion API
type ChatResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Agent defines the interface for interacting with an AI agent
type Agent interface {
	// SendMessage sends a message to the agent and returns the response
	SendMessage(ctx context.Context, prompt string) (string, error)
	
	// SendMessages sends multiple messages to the agent
	SendMessages(ctx context.Context, messages []Message) (ChatResponse, error)
	
	// GetStatus returns the current agent status
	GetStatus() Status
	
	// Connect establishes connection to the agent
	Connect(ctx context.Context) error
	
	// Disconnect closes connection to the agent
	Disconnect() error
}

// Status represents the agent's current status
type Status struct {
	Connected   bool   `json:"connected"`
	Endpoint    string `json:"endpoint"`
	Model       string `json:"model"`
	LastError   string `json:"last_error,omitempty"`
	TotalCalls  int    `json:"total_calls"`
	FailedCalls int    `json:"failed_calls"`
}

// NewAgent creates a new agent instance from configuration
func NewAgent(cfg config.Config) Agent {
	return &HTTPAgent{
		endpoint:    cfg.Agent.Endpoint,
		timeout:     cfg.Agent.Timeout.ToDuration(),
		model:       cfg.Agent.Model,
		temperature: cfg.Agent.Temperature,
		maxTokens:   cfg.Agent.MaxTokens,
	}
}
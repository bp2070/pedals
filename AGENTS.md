# Agent Package

The agent package provides the core interface for interacting with AI agents via HTTP.

## Overview

The package implements an OpenAI-compatible Chat Completions API client with connection management, status tracking, and error handling.

## Architecture

```
agent/
├── agent.go      # Interface definitions and types
├── client.go     # HTTP implementation
└── agent_test.go # Unit tests
```

## Agent Interface

```go
type Agent interface {
    SendMessage(ctx context.Context, prompt string) (string, error)
    SendMessages(ctx context.Context, messages []Message) (ChatResponse, error)
    GetStatus() Status
    Connect(ctx context.Context) error
    Disconnect() error
}
```

## Types

### Message
Represents a chat message with role and content.

```go
type Message struct {
    Role    string `json:"role"`    // "user", "assistant", "system"
    Content string `json:"content"`
}
```

### ChatRequest
Request payload sent to the agent endpoint.

```go
type ChatRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    Temperature float64   `json:"temperature,omitempty"`
    MaxTokens   int       `json:"max_tokens,omitempty"`
    Stream      bool      `json:"stream,omitempty"`
}
```

### ChatResponse
Response payload from the agent endpoint.

```go
type ChatResponse struct {
    ID      string   `json:"id"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}
```

### Status
Connection and usage statistics.

```go
type Status struct {
    Connected   bool   `json:"connected"`
    Endpoint    string `json:"endpoint"`
    Model       string `json:"model"`
    LastError   string `json:"last_error,omitempty"`
    TotalCalls  int    `json:"total_calls"`
    FailedCalls int    `json:"failed_calls"`
}
```

## HTTPAgent Implementation

`HTTPAgent` is the default implementation using standard library HTTP client.

### Configuration

```go
cfg := config.Config{
    Agent: config.AgentConfig{
        Endpoint:    "http://localhost:8080/chat/completions",
        Timeout:     config.Duration(30 * time.Second),
        Model:       "gpt-4",
        Temperature: 0.7,
        MaxTokens:   1000,
    },
}

ag := agent.NewAgent(cfg)
```

### Usage

```go
ctx := context.Background()

// Connect to agent
if err := ag.Connect(ctx); err != nil {
    log.Fatal(err)
}

// Send single message
response, err := ag.SendMessage(ctx, "Hello")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response)

// Send multiple messages (maintains conversation)
messages := []agent.Message{
    {Role: "system", Content: "You are a helpful assistant."},
    {Role: "user", Content: "Hello"},
}
resp, err := ag.SendMessages(ctx, messages)

// Check status
status := ag.GetStatus()
fmt.Printf("Connected: %v, Calls: %d/%d\n", status.Connected, status.TotalCalls, status.FailedCalls)

// Disconnect
ag.Disconnect()
```

## API Compatibility

The package expects OpenAI-compatible `/chat/completions` endpoints.

### Request
```json
{
  "model": "model-name",
  "messages": [{"role": "user", "content": "Hello"}],
  "temperature": 0.7,
  "max_tokens": 1000,
  "stream": false
}
```

### Response
```json
{
  "id": "chat-123",
  "model": "model-name",
  "choices": [{
    "index": 0,
    "message": {"role": "assistant", "content": "Response"},
    "finish_reason": "stop"
  }],
  "usage": {"prompt_tokens": 5, "completion_tokens": 10, "total_tokens": 15}
}
```

## Error Handling

The agent tracks errors in `Status.LastError` and maintains call statistics:

- `TotalCalls`: Successful API calls
- `FailedCalls`: Failed API calls (network errors, non-200 status, empty choices)

```go
status := ag.GetStatus()
if status.FailedCalls > 0 {
    fmt.Println("Last error:", status.LastError)
}
```

## Timeout

Requests are subject to two timeouts:
1. Context timeout passed to `SendMessage`/`SendMessages`
2. Client timeout configured via `config.Agent.Timeout`

The shorter of the two takes precedence.

## Extending

To implement a custom agent (e.g., WebSocket, gRPC), implement the `Agent` interface:

```go
type MyAgent struct {
    // custom fields
}

func (m *MyAgent) SendMessage(ctx context.Context, prompt string) (string, error) {
    // implementation
}

func (m *MyAgent) SendMessages(ctx context.Context, messages []Message) (ChatResponse, error) {
    // implementation
}

func (m *MyAgent) GetStatus() Status {
    // implementation
}

func (m *MyAgent) Connect(ctx context.Context) error {
    // implementation
}

func (m *MyAgent) Disconnect() error {
    // implementation
}
```
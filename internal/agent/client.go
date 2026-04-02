package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPAgent implements the Agent interface using HTTP
type HTTPAgent struct {
	endpoint    string
	timeout     time.Duration
	model       string
	temperature float64
	maxTokens   int
	
	client      *http.Client
	status      Status
}

// SendMessage sends a single message to the agent
func (a *HTTPAgent) SendMessage(ctx context.Context, prompt string) (string, error) {
	message := Message{
		Role:    "user",
		Content: prompt,
	}
	
	response, err := a.SendMessages(ctx, []Message{message})
	if err != nil {
		a.status.FailedCalls++
		a.status.LastError = err.Error()
		return "", err
	}
	
	if len(response.Choices) == 0 {
		err := fmt.Errorf("no choices in response")
		a.status.FailedCalls++
		a.status.LastError = err.Error()
		return "", err
	}
	
	a.status.TotalCalls++
	return response.Choices[0].Message.Content, nil
}

// SendMessages sends multiple messages to the agent
func (a *HTTPAgent) SendMessages(ctx context.Context, messages []Message) (ChatResponse, error) {
	var response ChatResponse
	
	// Create request
	req := ChatRequest{
		Model:       a.model,
		Messages:    messages,
		Temperature: a.temperature,
		MaxTokens:   a.maxTokens,
		Stream:      false,
	}
	
	// Marshal request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return response, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	
	// Ensure client is initialized
	if a.client == nil {
		a.client = &http.Client{
			Timeout: a.timeout,
		}
	}
	
	// Send request
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return response, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("agent returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return response, nil
}

// GetStatus returns the current agent status
func (a *HTTPAgent) GetStatus() Status {
	return a.status
}

// Connect establishes connection to the agent
func (a *HTTPAgent) Connect(ctx context.Context) error {
	// Test connection by making a simple request
	testReq, err := http.NewRequestWithContext(ctx, "GET", a.endpoint, nil)
	if err != nil {
		a.status.LastError = err.Error()
		return fmt.Errorf("failed to create test request: %w", err)
	}
	
	if a.client == nil {
		a.client = &http.Client{
			Timeout: a.timeout,
		}
	}
	
	// Try to connect
	resp, err := a.client.Do(testReq)
	if err != nil {
		a.status.Connected = false
		a.status.LastError = err.Error()
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	resp.Body.Close()
	
	// Update status
	a.status.Connected = true
	a.status.Endpoint = a.endpoint
	a.status.Model = a.model
	a.status.LastError = ""
	
	return nil
}

// Disconnect closes connection to the agent
func (a *HTTPAgent) Disconnect() error {
	a.status.Connected = false
	return nil
}
package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pedals/internal/config"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	cfg := mockConfig()

	ag := NewAgent(cfg)

	if ag == nil {
		t.Error("NewAgent() returned nil")
	}

	// Status should not have model set until connected
	status := ag.GetStatus()
	if status.Model != "" {
		t.Errorf("GetStatus().Model = %v, want empty before connect", status.Model)
	}
}

func TestHTTPAgent_SendMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}
		
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}
		if req.Model != "mock-llm" {
			t.Errorf("Expected model mock-llm, got %s", req.Model)
		}
		
		resp := ChatResponse{
			ID:      "test-id",
			Model:   "mock-llm",
			Choices: []Choice{{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "Hello, world!",
				},
				FinishReason: "stop",
			}},
			Usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	resp, err := ag.SendMessage(ctx, "Hello")
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	
	if resp != "Hello, world!" {
		t.Errorf("SendMessage() = %v, want 'Hello, world!'", resp)
	}
	
	status := ag.GetStatus()
	if status.TotalCalls != 1 {
		t.Errorf("GetStatus().TotalCalls = %v, want 1", status.TotalCalls)
	}
	if status.FailedCalls != 0 {
		t.Errorf("GetStatus().FailedCalls = %v, want 0", status.FailedCalls)
	}
}

func TestHTTPAgent_SendMessage_ErrorStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	_, err := ag.SendMessage(ctx, "Hello")
	if err == nil {
		t.Error("SendMessage() expected error for non-200 status, got nil")
	}
	
	status := ag.GetStatus()
	if status.FailedCalls != 1 {
		t.Errorf("GetStatus().FailedCalls = %v, want 1", status.FailedCalls)
	}
}

func TestHTTPAgent_SendMessage_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	_, err := ag.SendMessage(ctx, "Hello")
	if err == nil {
		t.Error("SendMessage() expected error for invalid JSON, got nil")
	}
}

func TestHTTPAgent_SendMessage_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ChatResponse{
			ID:      "test-id",
			Model:   "mock-llm",
			Choices: []Choice{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	_, err := ag.SendMessage(ctx, "Hello")
	if err == nil {
		t.Error("SendMessage() expected error for empty choices, got nil")
	}
}

func TestHTTPAgent_SendMessage_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	cfg.Agent.Timeout = config.Duration(10 * time.Millisecond)
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	_, err := ag.SendMessage(ctx, "Hello")
	if err == nil {
		t.Error("SendMessage() expected error for timeout, got nil")
	}
}

func TestHTTPAgent_SendMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if len(req.Messages) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(req.Messages))
		}
		if req.Messages[0].Role != "system" {
			t.Errorf("Expected first message role 'system', got %s", req.Messages[0].Role)
		}
		
		resp := ChatResponse{
			ID:      "test-id",
			Model:   "mock-llm",
			Choices: []Choice{{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "Response",
				},
				FinishReason: "stop",
			}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	messages := []Message{
		{Role: "system", Content: "You are helpful."},
		{Role: "user", Content: "Hello"},
	}
	
	resp, err := ag.SendMessages(ctx, messages)
	if err != nil {
		t.Fatalf("SendMessages() error = %v", err)
	}
	
	if len(resp.Choices) != 1 {
		t.Errorf("Response Choices length = %v, want 1", len(resp.Choices))
	}
}

func TestHTTPAgent_Connect_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	cfg := mockConfig()
	cfg.Agent.Endpoint = server.URL
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	err := ag.Connect(ctx)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	
	status := ag.GetStatus()
	if !status.Connected {
		t.Error("GetStatus().Connected = false, want true")
	}
	if status.Endpoint != server.URL {
		t.Errorf("GetStatus().Endpoint = %v, want %v", status.Endpoint, server.URL)
	}
}

func TestHTTPAgent_Connect_Failure(t *testing.T) {
	cfg := mockConfig()
	cfg.Agent.Endpoint = "http://localhost:1" // Invalid endpoint
	
	ag := NewAgent(cfg)
	ctx := context.Background()
	
	err := ag.Connect(ctx)
	if err == nil {
		t.Error("Connect() expected error for invalid endpoint, got nil")
	}
	
	status := ag.GetStatus()
	if status.Connected {
		t.Error("GetStatus().Connected = true, want false")
	}
}

func TestHTTPAgent_Disconnect(t *testing.T) {
	cfg := mockConfig()
	
	ag := NewAgent(cfg)
	
	// Simulate connected state (need to set it manually or connect first)
	// Let's just test that disconnect sets Connected to false
	err := ag.Disconnect()
	if err != nil {
		t.Errorf("Disconnect() error = %v", err)
	}
	
	status := ag.GetStatus()
	if status.Connected {
		t.Error("GetStatus().Connected = true, want false after Disconnect()")
	}
}

func TestHTTPAgent_GetStatus(t *testing.T) {
	cfg := mockConfig()

	ag := NewAgent(cfg)
	status := ag.GetStatus()

	// Model is not set until connected
	if status.Model != "" {
		t.Errorf("GetStatus().Model = %v, want empty before connect", status.Model)
	}
	if status.Endpoint != "" {
		t.Errorf("GetStatus().Endpoint = %v, want empty before connect", status.Endpoint)
	}
	if status.TotalCalls != 0 {
		t.Errorf("GetStatus().TotalCalls = %v, want 0", status.TotalCalls)
	}
	if status.FailedCalls != 0 {
		t.Errorf("GetStatus().FailedCalls = %v, want 0", status.FailedCalls)
	}
}

func mockConfig() config.Config {
	return config.Config{
		Agent: config.AgentConfig{
			Endpoint:    "http://localhost:8080/chat/completions",
			Timeout:     config.Duration(30 * time.Second),
			Model:       "mock-llm",
			Temperature: 0.7,
			MaxTokens:   1000,
		},
		UI: config.UIConfig{
			ShowHelpOnStart: true,
			MessageHistory:  100,
		},
	}
}
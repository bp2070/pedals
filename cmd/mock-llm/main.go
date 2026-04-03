package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"pedals/internal/agent"
	"time"
)

var (
	port     = flag.Int("port", 8090, "Port to listen on")
	delay    = flag.Int("delay", 0, "Artificial delay in milliseconds")
	behavior = flag.String("behavior", "echo", "Response behavior: echo, fixed, random")
	logReqs  = flag.Bool("log", true, "Log requests")
)

func main() {
	flag.Parse()

	http.HandleFunc("/chat/completions", chatHandler)
	http.HandleFunc("/health", healthHandler)
	
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Mock LLM server starting on %s", addr)
	log.Printf("Behavior: %s, Delay: %dms", *behavior, *delay)
	log.Printf("Health check: http://localhost:%d/health", *port)
	log.Printf("Chat endpoint: http://localhost:%d/chat/completions", *port)
	
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse request
	var req agent.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Log request if enabled
	if *logReqs {
		log.Printf("[REQUEST] Model: %s, Messages: %d, Temp: %.2f", 
			req.Model, len(req.Messages), req.Temperature)
		for i, msg := range req.Messages {
			content := msg.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			log.Printf("  [%d] %s: %s", i, msg.Role, content)
		}
	}

	// Artificial delay
	if *delay > 0 {
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}

	// Generate response based on behavior
	response := generateResponse(req)

	// Calculate token usage (rough approximation)
	promptTokens := estimateTokens(req.Messages)
	completionTokens := estimateTokens([]agent.Message{{Content: response}})

	// Build response
	resp := agent.ChatResponse{
		ID:    fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		Model: req.Model,
		Choices: []agent.Choice{
			{
				Index: 0,
				Message: agent.Message{
					Role:    "assistant",
					Content: response,
				},
				FinishReason: "stop",
			},
		},
		Usage: agent.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}
	
	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	
	// Log response time
	if *logReqs {
		log.Printf("[RESPONSE] Duration: %v, Tokens: %d/%d", 
			time.Since(start), completionTokens, resp.Usage.TotalTokens)
	}
}

func generateResponse(req agent.ChatRequest) string {
	if len(req.Messages) == 0 {
		return "I received an empty message. How can I help you?"
	}
	
	lastMessage := req.Messages[len(req.Messages)-1].Content
	
	switch *behavior {
	case "fixed":
		return "This is a fixed response from the mock LLM. Your message was: " + lastMessage
	case "random":
		responses := []string{
			"I understand you said: " + lastMessage,
			"Interesting point about: " + lastMessage,
			"Based on your input, I would recommend further analysis.",
			"Mock LLM response: Processing complete.",
			"I'm a mock LLM, but I'll pretend to think about: " + lastMessage,
		}
		return responses[rand.Intn(len(responses))]
	case "echo":
		fallthrough
	default:
		return "Mock LLM response to: " + lastMessage
	}
}

func estimateTokens(messages []agent.Message) int {
	// Rough approximation: 1 token ≈ 4 characters for English
	totalChars := 0
	for _, msg := range messages {
		totalChars += len(msg.Content)
	}
	return (totalChars + 3) / 4 // Round up
}
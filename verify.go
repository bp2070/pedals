// +build ignore

// This file is used to verify that the code compiles without errors.
// Run with: go run verify.go

package main

import (
	"fmt"
	"pedals/internal/agent"
	"pedals/internal/config"
	"pedals/internal/tui"
)

func main() {
	fmt.Println("Verifying project structure...")
	
	// Test config loading
	cfg := config.DefaultConfig()
	fmt.Printf("Default config loaded: %+v\n", cfg)
	
	// Test agent creation
	ag := agent.NewAgent(cfg)
	fmt.Printf("Agent created: %T\n", ag)
	
	// Test TUI model creation
	model := tui.NewModel(cfg, ag)
	fmt.Printf("TUI model created: %T\n", model)
	
	fmt.Println("Verification successful!")
}
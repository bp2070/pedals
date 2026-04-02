package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Agent AgentConfig `json:"agent"`
	UI    UIConfig    `json:"ui"`
}

// Duration is a wrapper for time.Duration that supports JSON marshaling/unmarshaling
type Duration time.Duration

// MarshalJSON implements json.Marshaler
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

// ToDuration returns the Duration as a time.Duration
func (d Duration) ToDuration() time.Duration {
	return time.Duration(d)
}

// AgentConfig contains agent-specific configuration
type AgentConfig struct {
	Endpoint    string   `json:"endpoint"`
	Timeout     Duration `json:"timeout"`
	Model       string   `json:"model"`
	Temperature float64  `json:"temperature"`
	MaxTokens   int      `json:"max_tokens"`
}

// UIConfig contains UI-specific configuration
type UIConfig struct {
	ShowHelpOnStart bool `json:"show_help_on_start"`
	MessageHistory  int  `json:"message_history"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Agent: AgentConfig{
			Endpoint:    "http://localhost:8080/chat/completions",
			Timeout:     Duration(30 * time.Second),
			Model:       "mock-llm",
			Temperature: 0.7,
			MaxTokens:   1000,
		},
		UI: UIConfig{
			ShowHelpOnStart: true,
			MessageHistory:  100,
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (Config, error) {
	config := DefaultConfig()

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Config file doesn't exist, return defaults
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON config
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
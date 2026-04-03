package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		want     string
		wantErr  bool
	}{
		{
			name:     "seconds",
			duration: Duration(30 * time.Second),
			want:     `"30s"`,
			wantErr:  false,
		},
		{
			name:     "minutes",
			duration: Duration(5 * time.Minute),
			want:     `"5m0s"`,
			wantErr:  false,
		},
		{
			name:     "zero",
			duration: Duration(0),
			want:     `"0s"`,
			wantErr:  false,
		},
		{
			name:     "complex",
			duration: Duration(1*time.Hour + 30*time.Minute + 15*time.Second),
			want:     `"1h30m15s"`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.duration.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Duration.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		want     Duration
		wantErr  bool
	}{
		{
			name:     "seconds",
			json:     `"30s"`,
			want:     Duration(30 * time.Second),
			wantErr:  false,
		},
		{
			name:     "minutes",
			json:     `"5m"`,
			want:     Duration(5 * time.Minute),
			wantErr:  false,
		},
		{
			name:     "zero",
			json:     `"0s"`,
			want:     Duration(0),
			wantErr:  false,
		},
		{
			name:     "invalid",
			json:     `"invalid"`,
			want:     Duration(0),
			wantErr:  true,
		},
		{
			name:     "not a string",
			json:     `123`,
			want:     Duration(0),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			err := d.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && d != tt.want {
				t.Errorf("Duration.UnmarshalJSON() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDuration_ToDuration(t *testing.T) {
	d := Duration(30 * time.Second)
	want := 30 * time.Second

	got := d.ToDuration()
	if got != want {
		t.Errorf("Duration.ToDuration() = %v, want %v", got, want)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Agent.Endpoint != "http://localhost:8090/chat/completions" {
		t.Errorf("DefaultConfig().Agent.Endpoint = %v, want http://localhost:8090/chat/completions", cfg.Agent.Endpoint)
	}
	if cfg.Agent.Timeout.ToDuration() != 30*time.Second {
		t.Errorf("DefaultConfig().Agent.Timeout = %v, want 30s", cfg.Agent.Timeout)
	}
	if cfg.Agent.Model != "mock-llm" {
		t.Errorf("DefaultConfig().Agent.Model = %v, want mock-llm", cfg.Agent.Model)
	}
	if cfg.Agent.Temperature != 0.7 {
		t.Errorf("DefaultConfig().Agent.Temperature = %v, want 0.7", cfg.Agent.Temperature)
	}
	if cfg.Agent.MaxTokens != 1000 {
		t.Errorf("DefaultConfig().Agent.MaxTokens = %v, want 1000", cfg.Agent.MaxTokens)
	}
	if !cfg.UI.ShowHelpOnStart {
		t.Errorf("DefaultConfig().UI.ShowHelpOnStart = %v, want true", cfg.UI.ShowHelpOnStart)
	}
	if cfg.UI.MessageHistory != 100 {
		t.Errorf("DefaultConfig().UI.MessageHistory = %v, want 100", cfg.UI.MessageHistory)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/path/config.json")

	if err != nil {
		t.Errorf("LoadConfig() error = %v, want nil", err)
	}
	if cfg.Agent.Endpoint != "http://localhost:8090/chat/completions" {
		t.Errorf("LoadConfig() returned non-default config when file doesn't exist")
	}
}

func TestLoadConfig_SaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	original := Config{
		Agent: AgentConfig{
			Endpoint:    "http://custom:9999/chat",
			Timeout:     Duration(45 * time.Second),
			Model:       "gpt-4",
			Temperature: 0.9,
			MaxTokens:   2000,
		},
		UI: UIConfig{
			ShowHelpOnStart: false,
			MessageHistory:  50,
		},
	}

	err := SaveConfig(configPath, original)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if loaded.Agent.Endpoint != original.Agent.Endpoint {
		t.Errorf("Agent.Endpoint = %v, want %v", loaded.Agent.Endpoint, original.Agent.Endpoint)
	}
	if loaded.Agent.Timeout != original.Agent.Timeout {
		t.Errorf("Agent.Timeout = %v, want %v", loaded.Agent.Timeout, original.Agent.Timeout)
	}
	if loaded.Agent.Model != original.Agent.Model {
		t.Errorf("Agent.Model = %v, want %v", loaded.Agent.Model, original.Agent.Model)
	}
	if loaded.Agent.Temperature != original.Agent.Temperature {
		t.Errorf("Agent.Temperature = %v, want %v", loaded.Agent.Temperature, original.Agent.Temperature)
	}
	if loaded.Agent.MaxTokens != original.Agent.MaxTokens {
		t.Errorf("Agent.MaxTokens = %v, want %v", loaded.Agent.MaxTokens, original.Agent.MaxTokens)
	}
	if loaded.UI.ShowHelpOnStart != original.UI.ShowHelpOnStart {
		t.Errorf("UI.ShowHelpOnStart = %v, want %v", loaded.UI.ShowHelpOnStart, original.UI.ShowHelpOnStart)
	}
	if loaded.UI.MessageHistory != original.UI.MessageHistory {
		t.Errorf("UI.MessageHistory = %v, want %v", loaded.UI.MessageHistory, original.UI.MessageHistory)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(configPath, []byte(`{invalid json`), 0644); err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid JSON, got nil")
	}
}

func TestSaveConfig_FileWriteError(t *testing.T) {
	cfg := DefaultConfig()
	
	// Try to save to a directory that doesn't exist
	err := SaveConfig("/nonexistent directory/config.json", cfg)
	if err == nil {
		t.Error("SaveConfig() expected error for unwritable path, got nil")
	}
}
# Pedals - Minimalistic TUI Agent Harness

A lightweight Text User Interface (TUI) for interacting with custom AI agents using OpenAI-like APIs, built with Go and Bubble Tea.

## Features

- **Clean TUI Interface**: Chat interface with agent responses
- **Agent Control**: Connect/disconnect, send messages, view status
- **OpenAI-like API**: Compatible with `/chat/completions` endpoints
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Configuration**: JSON-based config with sensible defaults

## Installation

### Prerequisites
- Go 1.21 or later

### Build from Source
```bash
# Clone the repository
git clone <repository-url>
cd pedals

# Build the binary
go build -o pedals ./cmd/pedals

# Run the application
./pedals
```

### Configuration
Create a `config.json` file in the same directory as the binary:

```json
{
  "agent": {
    "endpoint": "http://localhost:8080/chat/completions",
    "timeout": "30s",
    "model": "custom-model",
    "temperature": 0.7,
    "max_tokens": 1000
  },
  "ui": {
    "show_help_on_start": true,
    "message_history": 100
  }
}
```

## Usage

### Basic Usage
1. Configure your agent endpoint in `config.json`
2. Run the application: `./pedals`
3. Type messages and press Enter to send
4. View agent responses in the chat window

### Keybindings
- **Enter**: Send message
- **Ctrl+R**: Reconnect to agent
- **Ctrl+L**: Clear chat history
- **Ctrl+C** or **Esc**: Quit application

### Status Bar
The status bar at the bottom shows:
- Connection status (Connected/Disconnected)
- Agent endpoint
- Model name
- Call statistics (successful/failed)
- Loading indicator when processing

## Project Structure

```
pedals/
├── cmd/pedals/main.go          # Application entry point
├── internal/
│   ├── agent/                  # Agent control logic
│   │   ├── agent.go           # Agent interface and types
│   │   └── client.go          # HTTP client implementation
│   ├── tui/                   # TUI components
│   │   ├── model.go           # Bubble Tea model
│   │   ├── styles.go          # UI styling
│   │   └── views/             # UI components (future)
│   └── config/                # Configuration
│       └── config.go          # Config loading/parsing
├── config.json                # Example configuration
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Development

### Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea): TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): Styling
- [Bubbles](https://github.com/charmbracelet/bubbles): UI components

### Adding Features
1. Agent interface is defined in `internal/agent/agent.go`
2. TUI model is in `internal/tui/model.go`
3. Configuration is handled in `internal/config/config.go`

### Testing
```bash
# Run tests
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o pedals-linux ./cmd/pedals
GOOS=darwin GOARCH=arm64 go build -o pedals-macos ./cmd/pedals
GOOS=windows GOARCH=amd64 go build -o pedals.exe ./cmd/pedals
```

## API Compatibility

The agent harness expects endpoints that follow the OpenAI Chat Completions API format:

### Request Format
```json
{
  "model": "custom-model",
  "messages": [
    {"role": "user", "content": "Hello"}
  ],
  "temperature": 0.7,
  "max_tokens": 1000,
  "stream": false
}
```

### Response Format
```json
{
  "id": "chat-123",
  "model": "custom-model",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello there!"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 5,
    "completion_tokens": 7,
    "total_tokens": 12
  }
}
```

## License

MIT

## Acknowledgments

- [Charmbracelet](https://charmbracelet.com/) for the excellent Bubble Tea framework
- Inspired by various TUI tools and agent frameworks
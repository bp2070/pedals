# Minimalistic TUI Agent Harness - Implementation Plan

## Overview
A Go-based TUI (Text User Interface) agent harness using the Bubble Tea framework for controlling custom AI agents with OpenAI-like APIs.

## Requirements Summary
- **Language**: Go with Bubble Tea TUI framework
- **Target**: Custom AI agent with OpenAI-like `/chat/completions` API
- **MVP**: Basic TUI with agent control
- **Platform**: Cross-platform (Windows/Linux/macOS)
- **Deployment**: Local development only

## Architecture

### Core Components
1. **TUI Framework**: Bubble Tea (model-update-view pattern)
2. **HTTP Client**: Standard library or lightweight HTTP client for API calls
3. **Configuration**: Simple config file (JSON/YAML) for agent endpoints
4. **State Management**: Bubble Tea's Model for managing application state

### Project Structure
```
pedals/
├── cmd/
│   └── pedals/
│       └── main.go          # Application entry point
├── internal/
│   ├── agent/               # Agent control logic
│   │   ├── agent.go         # Agent interface and implementation
│   │   └── client.go        # HTTP client for API calls
│   ├── tui/                 # TUI components
│   │   ├── model.go         # Bubble Tea Model
│   │   ├── views/           # UI components
│   │   │   ├── chat.go      # Chat interface
│   │   │   ├── status.go    # Status display
│   │   │   └── input.go     # Command input
│   │   └── update.go        # Message handlers
│   └── config/              # Configuration
│       └── config.go        # Config loading/parsing
├── config.yaml              # Example configuration
├── go.mod                   # Go module definition
└── README.md                # Project documentation
```

## Implementation Phases

### Phase 1: Project Setup
1. Initialize Go module
2. Install dependencies (Bubble Tea, Lip Gloss, Bubbles)
3. Create basic directory structure

### Phase 2: Core TUI Structure
1. Implement Bubble Tea Model with basic state
2. Create minimal view rendering
3. Set up message passing system
4. Add quit/help keybindings

### Phase 3: Agent Integration
1. Create HTTP client for OpenAI-like API
2. Implement agent interface with start/stop methods
3. Add basic agent status tracking
4. Handle API errors and timeouts

### Phase 4: Chat Interface
1. Implement chat message display
2. Add text input for user prompts
3. Create message history storage
4. Add send/clear commands

### Phase 5: Configuration & Polish
1. Add config file support (YAML/JSON)
2. Implement help screen
3. Add status indicators
4. Improve error handling

### Phase 6: Testing & Documentation
1. Add example configuration
2. Create README with usage instructions
3. Test cross-platform compatibility
4. Add build scripts

## Dependencies

### Required Packages
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components (textinput, spinner, etc.)
- `gopkg.in/yaml.v3` or similar for config parsing

### Optional/Recommended
- `github.com/spf13/cobra` - CLI structure (if adding CLI options)
- `github.com/spf13/viper` - Config management (if complex config needed)

## MVP Features

### TUI Interface
1. **Status Bar**: Agent connection status, API endpoint, basic stats
2. **Chat Area**: Display conversation history
3. **Input Area**: Text input for user prompts
4. **Help Screen**: Keybindings reference (F1 or ?)
5. **Quit**: Ctrl+C or Ctrl+Q to exit

### Agent Control
1. **Connection Management**: Connect/disconnect to agent API
2. **Chat Completions**: Send prompts, receive responses
3. **Status Monitoring**: Track request/response counts, errors
4. **Basic Configuration**: Set API endpoint, timeouts

### Configuration
1. **Endpoint URL**: Agent API endpoint
2. **API Key**: Authentication (if needed)
3. **Request Timeout**: Timeout for API calls
4. **Model/Parameters**: Default model, temperature, etc.

## Key Design Decisions

### 1. State Management
- Use Bubble Tea's Model pattern
- Keep state simple: agent status, chat history, input buffer
- Use messages for all state transitions

### 2. API Client Design
- Simple interface: `SendMessage(prompt string) (response string, err error)`
- Support streaming responses if needed later
- Basic error handling with retry logic

### 3. Configuration
- YAML format for readability
- Default config with overrides from file
- Environment variable support for sensitive data

### 4. Cross-Platform Considerations
- Use Bubble Tea's built-in cross-platform support
- Avoid OS-specific file paths or commands
- Test line ending handling for Windows

## Success Criteria
1. ✅ TUI starts without errors
2. ✅ Can connect to agent API endpoint
3. ✅ Send prompts and receive responses
4. ✅ Display conversation history
5. ✅ Basic agent status monitoring
6. ✅ Clean exit with Ctrl+C/Ctrl+Q
7. ✅ Cross-platform compatibility (macOS/Linux/Windows)

## Next Steps After MVP
1. Add streaming response display
2. Implement conversation persistence
3. Add multiple agent support
4. Include advanced configuration options
5. Add metrics and logging
6. Package as single binary with releases

## Timeline Estimate
- **Phase 1-2**: 1-2 hours (basic TUI skeleton)
- **Phase 3-4**: 2-3 hours (agent integration + chat)
- **Phase 5-6**: 1-2 hours (polish + documentation)

Total: 4-7 hours for complete MVP
package tui

import (
	"context"
	"fmt"
	"pedals/internal/agent"
	"pedals/internal/config"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI application state
type Model struct {
	// Configuration
	config config.Config
	
	// Agent
	agent agent.Agent
	
	// UI components
	viewport    viewport.Model
	textInput   textinput.Model
	spinner     spinner.Model
	
	// State
	ready       bool
	width       int
	height      int
	quitting    bool
	loading     bool
	messages    []agent.Message
	status      agent.Status
	
	// Styles
	styles *Styles
}

// NewModel creates a new TUI model
func NewModel(cfg config.Config, ag agent.Agent) Model {
	// Create text input
	ti := textinput.New()
	ti.Placeholder = "Type your message here..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 50
	
	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	
	// Create viewport
	vp := viewport.New(0, 0)
	
	// Create styles
	styles := DefaultStyles()
	
	return Model{
		config:      cfg,
		agent:       ag,
		textInput:   ti,
		spinner:     s,
		viewport:    vp,
		ready:       false,
		quitting:    false,
		loading:     false,
		messages:    make([]agent.Message, 0),
		styles:      styles,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Start with spinner and connection attempt
	return tea.Batch(
		m.spinner.Tick,
		connectToAgent(m.agent),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		if !m.ready {
			m.ready = true
		}
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4 // Reserve space for input and status
		m.textInput.Width = msg.Width - 10
		m.updateViewportContent()
		
	case tea.KeyMsg:
		// Handle key presses
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if !m.loading && m.textInput.Value() != "" {
				// Send message
				message := m.textInput.Value()
				m.textInput.SetValue("")
				m.messages = append(m.messages, agent.Message{
					Role:    "user",
					Content: message,
				})
				m.updateViewportContent()
				m.loading = true
				return m, tea.Batch(
					m.spinner.Tick,
					sendMessage(m.agent, message),
				)
			}
		case "ctrl+r":
			// Reconnect
			m.loading = true
			return m, tea.Batch(
				m.spinner.Tick,
				connectToAgent(m.agent),
			)
		case "ctrl+l":
			// Clear messages
			m.messages = make([]agent.Message, 0)
			m.updateViewportContent()
		}
		
	case ConnectionResultMsg:
		// Handle connection result
		m.loading = false
		m.status = m.agent.GetStatus()
		if msg.Error != nil {
			// Add error message to chat
			m.messages = append(m.messages, agent.Message{
				Role:    "system",
				Content: fmt.Sprintf("Connection error: %v", msg.Error),
			})
		} else {
			m.messages = append(m.messages, agent.Message{
				Role:    "system",
				Content: "Connected to agent",
			})
		}
		m.updateViewportContent()
		
	case MessageResultMsg:
		// Handle message response
		m.loading = false
		m.status = m.agent.GetStatus()
		if msg.Error != nil {
			// Add error message to chat
			m.messages = append(m.messages, agent.Message{
				Role:    "system",
				Content: fmt.Sprintf("Error: %v", msg.Error),
			})
		} else {
			// Add agent response to chat
			m.messages = append(m.messages, agent.Message{
				Role:    "assistant",
				Content: msg.Response,
			})
		}
		m.updateViewportContent()
	}
	
	// Update components
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing...\n"
	}
	
	if m.quitting {
		return "Goodbye!\n"
	}
	
	// Build status bar
	status := m.buildStatusBar()
	
	// Build chat view
	chat := m.viewport.View()
	
	// Build input area
	input := m.buildInputArea()
	
	// Combine everything
	return lipgloss.JoinVertical(
		lipgloss.Left,
		chat,
		status,
		input,
	)
}

// updateViewportContent updates the viewport with current messages
func (m *Model) updateViewportContent() {
	var content strings.Builder
	
	for _, msg := range m.messages {
		switch msg.Role {
		case "user":
			content.WriteString(m.styles.UserStyle.Render("You: " + msg.Content))
		case "assistant":
			content.WriteString(m.styles.AgentStyle.Render("Agent: " + msg.Content))
		case "system":
			content.WriteString(m.styles.SystemStyle.Render("System: " + msg.Content))
		}
		content.WriteString("\n\n")
	}
	
	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

// buildStatusBar builds the status bar
func (m *Model) buildStatusBar() string {
	statusText := "Status: "
	if m.status.Connected {
		statusText += m.styles.ConnectedStyle.Render("Connected")
	} else {
		statusText += m.styles.DisconnectedStyle.Render("Disconnected")
	}
	
	statusText += fmt.Sprintf(" | Endpoint: %s", m.status.Endpoint)
	statusText += fmt.Sprintf(" | Model: %s", m.status.Model)
	statusText += fmt.Sprintf(" | Calls: %d/%d", m.status.TotalCalls, m.status.FailedCalls)
	
	if m.loading {
		statusText += " | " + m.spinner.View() + " Processing..."
	}
	
	return m.styles.StatusBarStyle.Render(statusText)
}

// buildInputArea builds the input area
func (m *Model) buildInputArea() string {
	input := m.textInput.View()
	help := "Enter: Send | Ctrl+R: Reconnect | Ctrl+L: Clear | Esc: Quit"
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.InputStyle.Render(input),
		m.styles.HelpStyle.Render(help),
	)
}

// Messages for tea.Cmd

// ConnectionResultMsg is sent when connection attempt completes
type ConnectionResultMsg struct {
	Error error
}

// MessageResultMsg is sent when message sending completes
type MessageResultMsg struct {
	Response string
	Error    error
}

// connectToAgent attempts to connect to the agent
func connectToAgent(ag agent.Agent) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		err := ag.Connect(ctx)
		return ConnectionResultMsg{Error: err}
	}
}

// sendMessage sends a message to the agent
func sendMessage(ag agent.Agent, message string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		response, err := ag.SendMessage(ctx, message)
		return MessageResultMsg{Response: response, Error: err}
	}
}
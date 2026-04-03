package tui

import "github.com/charmbracelet/lipgloss"

// Styles defines the UI styles for the application
type Styles struct {
	// General styles
	BorderStyle lipgloss.Style
	TitleStyle  lipgloss.Style
	
	// Status styles
	StatusBarStyle    lipgloss.Style
	ConnectedStyle    lipgloss.Style
	DisconnectedStyle lipgloss.Style
	
	// Message styles
	UserStyle    lipgloss.Style
	AgentStyle   lipgloss.Style
	SystemStyle  lipgloss.Style
	
	// Input styles
	InputStyle lipgloss.Style
	HelpStyle  lipgloss.Style
}

// DefaultStyles returns the default application styles
func DefaultStyles() *Styles {
	// Define colors
	primaryColor := lipgloss.Color("205")
	secondaryColor := lipgloss.Color("240")
	successColor := lipgloss.Color("46")
	errorColor := lipgloss.Color("160")
	userColor := lipgloss.Color("39")
	agentColor := lipgloss.Color("213")
	systemColor := lipgloss.Color("245")
	
	return &Styles{
		BorderStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1),
		
		TitleStyle: lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1),
		
		StatusBarStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(secondaryColor).
			Padding(0, 1),
		
		ConnectedStyle: lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true),
		
		DisconnectedStyle: lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true),
		
		UserStyle: lipgloss.NewStyle().
			Foreground(userColor).
			Bold(true),
		
		AgentStyle: lipgloss.NewStyle().
			Foreground(agentColor).
			Italic(true),
		
		SystemStyle: lipgloss.NewStyle().
			Foreground(systemColor).
			Faint(true),
		
		InputStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1),
		
		HelpStyle: lipgloss.NewStyle().
			Foreground(secondaryColor).
			Faint(true).
			Padding(0, 1),
	}
}
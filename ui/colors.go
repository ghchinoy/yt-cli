package ui

import "github.com/charmbracelet/lipgloss"

var (
	Accent  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // Blue
	Command = lipgloss.NewStyle().Foreground(lipgloss.Color("248")) // Light Grey
	Pass    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green
	Warn    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange
	Fail    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
	Muted   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Dark Grey
	ID      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // Teal/Mint
)

// FormatError provides a standard format for hints for agents
func FormatError(err error, hint string) string {
	out := Fail.Render("Error: ") + err.Error() + "\n"
	if hint != "" {
		out += Warn.Render("Hint: ") + Command.Render(hint) + "\n"
	}
	return out
}

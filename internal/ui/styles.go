package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Error     = lipgloss.Color("#EF4444") // Red
	Muted     = lipgloss.Color("#6B7280") // Gray

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	CodeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#F9FAFB")).
			Padding(0, 1)

	StepStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	StepNumStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(Primary).
			Padding(0, 1).
			MarginRight(1)
)

// Icons
const (
	IconSuccess = "✨"
	IconCheck   = "✅"
	IconWarning = "⚠️"
	IconError   = "❌"
	IconInfo    = "💡"
	IconFolder  = "📁"
	IconFile    = "📝"
	IconSync    = "🔄"
)

// Title prints the app title banner
func Title() string {
	banner := `
  ╭─────────────────────────────────────╮
  │  📝 readme-gen                      │
  ╰─────────────────────────────────────╯`
	return TitleStyle.Render(banner)
}

// Success prints a success message
func Success(msg string) string {
	return SuccessStyle.Render(fmt.Sprintf("%s %s", IconSuccess, msg))
}

// Check prints a check message
func Check(msg string) string {
	return SuccessStyle.Render(fmt.Sprintf("%s %s", IconCheck, msg))
}

// Warn prints a warning message
func Warn(msg string) string {
	return WarningStyle.Render(fmt.Sprintf("%s %s", IconWarning, msg))
}

// Err prints an error message
func Err(msg string) string {
	return ErrorStyle.Render(fmt.Sprintf("%s %s", IconError, msg))
}

// Info prints an info message
func Info(msg string) string {
	return MutedStyle.Render(fmt.Sprintf("%s %s", IconInfo, msg))
}

// Box wraps content in a styled box
func Box(content string) string {
	return BoxStyle.Render(content)
}

// Step prints a step indicator
func Step(num int, total int, label string) string {
	return fmt.Sprintf("%s %s",
		StepNumStyle.Render(fmt.Sprintf("%d/%d", num, total)),
		StepStyle.Render(label))
}

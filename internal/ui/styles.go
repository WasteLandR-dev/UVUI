// Package ui provides the user interface for the application.
package ui

import "github.com/charmbracelet/lipgloss"

// Theme colors
const (
	PrimaryColor   = "#7C3AED"
	SecondaryColor = "#1F2937"
	AccentColor    = "#10B981"
	WarningColor   = "#F59E0B"
	ErrorColor     = "#EF4444"
	TextColor      = "#FFFFFF"
	MutedTextColor = "#64748B"
	BorderColor    = "#374151"
	SuccessColor   = "#10B981"
)

var (
	// TitleStyle is the style for titles.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(PrimaryColor)).
			Background(lipgloss.Color(SecondaryColor)).
			Padding(0, 1)

	// ActiveTabStyle is the style for active tabs.
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(TextColor)).
			Background(lipgloss.Color(PrimaryColor)).
			Padding(0, 1)

	// InactiveTabStyle is the style for inactive tabs.
	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(MutedTextColor)).
				Background(lipgloss.Color(SecondaryColor)).
				Padding(0, 1)

	// PanelStyle is the style for panels.
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(BorderColor)).
			Padding(1).
			Margin(1, 0)

	// ActivePanelStyle is the style for active panels.
	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(PrimaryColor)).
				Padding(1).
				Margin(1, 0)

	// StatusStyle is the style for the status bar.
	StatusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(SecondaryColor)).
			Foreground(lipgloss.Color("#D1D5DB")).
			Padding(0, 1)

	// HelpStyle is the style for the help text.
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Padding(0, 1)

	// SelectedItemStyle is the style for selected items in a list.
	SelectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(TextColor)).
				Background(lipgloss.Color(PrimaryColor)).
				Padding(0, 1)

	// UnselectedItemStyle is the style for unselected items in a list.
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D5DB")).
				Padding(0, 1)

	// InstalledVersionStyle is the style for installed Python versions.
	InstalledVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor))

	// AvailableVersionStyle is the style for available Python versions.
	AvailableVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(MutedTextColor))

	// CurrentVersionStyle is the style for the current Python version.
	CurrentVersionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(AccentColor))

	// LoadingStyle is the style for loading indicators.
	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(WarningColor))

	// SuccessStyle is the style for success messages.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(SuccessColor))

	// ErrorStyle is the style for error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ErrorColor))

	// InfoMessageStyle is the style for info messages.
	InfoMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#60A5FA"))

	// WarningMessageStyle is the style for warning messages.
	WarningMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(WarningColor))

	// ErrorMessageStyle is the style for error messages.
	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ErrorColor))

	// SuccessMessageStyle is the style for success messages.
	SuccessMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor))
)

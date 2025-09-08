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
	// TitleStyle defines the style for application titles.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(PrimaryColor)).
			Background(lipgloss.Color(SecondaryColor)).
			Padding(0, 1)

	// ActiveTabStyle defines the style for active navigation tabs.
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(TextColor)).
			Background(lipgloss.Color(PrimaryColor)).
			Padding(0, 1)

	// InactiveTabStyle defines the style for inactive navigation tabs.
	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(MutedTextColor)).
				Background(lipgloss.Color(SecondaryColor)).
				Padding(0, 1)

	// PanelStyle defines the default style for content panels.
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(BorderColor)).
			Padding(1).
			Margin(1, 0)

	// ActivePanelStyle defines the style for the currently active content panel.
	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(PrimaryColor)).
				Padding(1).
				Margin(1, 0)

	// StatusStyle defines the style for the application's status bar.
	StatusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(SecondaryColor)).
			Foreground(lipgloss.Color("#D1D5DB")).
			Padding(0, 1)

	// HelpStyle defines the style for help text.
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Padding(0, 1)

	// SelectedItemStyle defines the style for selected items in lists.
	SelectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(TextColor)).
				Background(lipgloss.Color(PrimaryColor)).
				Padding(0, 1)

	// UnselectedItemStyle defines the style for unselected items in lists.
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D5DB")).
				Padding(0, 1)

	// InstalledVersionStyle defines the style for installed Python versions.
	InstalledVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor))

	// AvailableVersionStyle defines the style for available Python versions.
	AvailableVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(MutedTextColor))

	// CurrentVersionStyle defines the style for the currently active Python version.
	CurrentVersionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(AccentColor))

	// LoadingStyle defines the style for loading indicators.
	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(WarningColor))

	// SuccessStyle defines the style for success messages.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(SuccessColor))

	// ErrorStyle defines the style for error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ErrorColor))

	// InfoMessageStyle defines the style for informational messages.
	InfoMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#60A5FA"))

	// WarningMessageStyle defines the style for warning messages.
	WarningMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(WarningColor))

	// ErrorMessageStyle defines the style for error messages.
	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ErrorColor))

	// SuccessMessageStyle defines the style for success messages.
	SuccessMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor))
)

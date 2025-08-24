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
	// Header styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(PrimaryColor)).
			Background(lipgloss.Color(SecondaryColor)).
			Padding(0, 1)

	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(TextColor)).
			Background(lipgloss.Color(PrimaryColor)).
			Padding(0, 1)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(MutedTextColor)).
				Background(lipgloss.Color(SecondaryColor)).
				Padding(0, 1)

	// Panel styles
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(BorderColor)).
			Padding(1).
			Margin(1, 0)

	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(PrimaryColor)).
				Padding(1).
				Margin(1, 0)

	// Status bar styles
	StatusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(SecondaryColor)).
			Foreground(lipgloss.Color("#D1D5DB")).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Padding(0, 1)

	// List item styles
	SelectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(TextColor)).
				Background(lipgloss.Color(PrimaryColor)).
				Padding(0, 1)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D5DB")).
				Padding(0, 1)

	// Version status styles
	InstalledVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor))

	AvailableVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(MutedTextColor))

	CurrentVersionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(AccentColor))

	// Status indicator styles
	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(WarningColor))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(SuccessColor))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ErrorColor))

	// Message styles
	InfoMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#60A5FA"))

	WarningMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(WarningColor))

	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ErrorColor))

	SuccessMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor))
)

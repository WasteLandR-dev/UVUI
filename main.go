package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#64748B")).
				Background(lipgloss.Color("#1F2937")).
				Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(1).
			Margin(1, 0)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7C3AED")).
				Padding(1).
				Margin(1, 0)

	statusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#D1D5DB")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Padding(0, 1)
)

// Panel represents a UI panel
type Panel int

const (
	StatusPanel Panel = iota
	PythonPanel
	ProjectPanel
	EnvironmentPanel
)

// UVStatus represents the status of UV installation
type UVStatus struct {
	Installed bool
	Version   string
	Path      string
}

// Model represents the application state
type Model struct {
	activePanel Panel
	panels      []Panel
	width       int
	height      int
	uvStatus    UVStatus
	installing  bool
	messages    []string
	uvInstaller UVInstallerInterface
}

// UVInstallerInterface defines the contract for UV installation
type UVInstallerInterface interface {
	IsInstalled() (bool, string, error)
	Install() error
	GetInstallCommand() (string, error)
}

// UVInstaller implements UV installation functionality
type UVInstaller struct{}

// NewUVInstaller creates a new UV installer
func NewUVInstaller() *UVInstaller {
	return &UVInstaller{}
}

// IsInstalled checks if UV is installed and returns version info
func (u *UVInstaller) IsInstalled() (bool, string, error) {
	// This would typically execute `uv --version` command
	// For now, returning mock data
	return false, "", nil
}

// Install installs UV based on the detected OS
func (u *UVInstaller) Install() error {
	// This would execute the appropriate installation command
	// Implementation depends on OS detection
	return nil
}

// GetInstallCommand returns the installation command for the current OS
func (u *UVInstaller) GetInstallCommand() (string, error) {
	// Detect OS and return appropriate command
	switch getOS() {
	case "windows":
		return "powershell -ExecutionPolicy ByPass -c \"irm https://astral.sh/uv/install.ps1 | iex\"", nil
	case "macos", "linux":
		return "curl -LsSf https://astral.sh/uv/install.sh | sh", nil
	default:
		return "", fmt.Errorf("unsupported operating system")
	}
}

// getOS detects the current operating system
func getOS() string {
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		return "windows"
	case fileExists("/System/Library/CoreServices/SystemVersion.plist"):
		return "macos"
	default:
		return "linux"
	}
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// InstallUVMsg represents a message to install UV
type InstallUVMsg struct{}

// UVInstalledMsg represents a message when UV installation is complete
type UVInstalledMsg struct {
	Success bool
	Version string
	Error   error
}

// NewModel creates a new application model
func NewModel() Model {
	return Model{
		activePanel: StatusPanel,
		panels: []Panel{
			StatusPanel,
			PythonPanel,
			ProjectPanel,
			EnvironmentPanel,
		},
		uvInstaller: NewUVInstaller(),
		messages:    []string{},
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		checkUVStatus(m.uvInstaller),
		tea.EnterAltScreen,
	)
}

// checkUVStatus checks the UV installation status
func checkUVStatus(installer UVInstallerInterface) tea.Cmd {
	return func() tea.Msg {
		installed, version, err := installer.IsInstalled()
		if err != nil {
			return UVInstalledMsg{Success: false, Version: version, Error: err}
		}
		return UVInstalledMsg{Success: installed, Version: version, Error: nil}
	}
}

// installUV installs UV
func installUV(installer UVInstallerInterface) tea.Cmd {
	return func() tea.Msg {
		err := installer.Install()
		return UVInstalledMsg{Success: err == nil, Version: "", Error: err}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			m.activePanel = Panel((int(m.activePanel) + 1) % len(m.panels))
			return m, nil

		case "shift+tab":
			if m.activePanel == 0 {
				m.activePanel = Panel(len(m.panels) - 1)
			} else {
				m.activePanel = Panel(int(m.activePanel) - 1)
			}
			return m, nil

		case "i":
			if !m.uvStatus.Installed && !m.installing {
				m.installing = true
				m.addMessage("Installing UV...")
				return m, installUV(m.uvInstaller)
			}
			return m, nil

		case "r":
			m.addMessage("Refreshing UV status...")
			return m, checkUVStatus(m.uvInstaller)

		case "?":
			m.addMessage("Help: q=quit, tab=next panel, shift+tab=prev panel, i=install UV, r=refresh")
			return m, nil
		}

	case UVInstalledMsg:
		m.installing = false
		if msg.Success {
			m.uvStatus.Installed = true
			m.addMessage("UV installation completed successfully!")
		} else if msg.Error != nil {
			m.addMessage(fmt.Sprintf("UV installation failed: %v", msg.Error))
		} else {
			m.uvStatus.Installed = false
			m.addMessage("UV is not installed")
		}
		return m, nil
	}

	return m, nil
}

// addMessage adds a message to the message list
func (m *Model) addMessage(msg string) {
	m.messages = append(m.messages, msg)
	// Keep only last 10 messages
	if len(m.messages) > 10 {
		m.messages = m.messages[1:]
	}
}

// View renders the application UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Header
	header := titleStyle.Render("UV Package Manager CLI")

	// Tabs
	tabs := m.renderTabs()

	// Main content based on active panel
	content := m.renderActivePanel()

	// Status bar
	statusBar := m.renderStatusBar()

	// Help line
	helpLine := helpStyle.Render("? for help | q to quit | tab to navigate")

	// Combine all elements
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabs,
		content,
		statusBar,
		helpLine,
	)
}

// renderTabs renders the navigation tabs
func (m Model) renderTabs() string {
	tabNames := []string{"Status", "Python", "Project", "Environment"}
	var tabs []string

	for i, name := range tabNames {
		if Panel(i) == m.activePanel {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderActivePanel renders the content of the active panel
func (m Model) renderActivePanel() string {
	style := panelStyle
	if m.activePanel == StatusPanel {
		style = activePanelStyle
	}

	switch m.activePanel {
	case StatusPanel:
		return style.Render(m.renderStatusPanel())
	case PythonPanel:
		return style.Render(m.renderPythonPanel())
	case ProjectPanel:
		return style.Render(m.renderProjectPanel())
	case EnvironmentPanel:
		return style.Render(m.renderEnvironmentPanel())
	default:
		return style.Render("Unknown panel")
	}
}

// renderStatusPanel renders the status panel content
func (m Model) renderStatusPanel() string {
	var content strings.Builder

	content.WriteString("UV Installation Status:\n\n")

	if m.installing {
		content.WriteString("⏳ Installing UV...\n")
	} else if m.uvStatus.Installed {
		content.WriteString("✅ UV is installed\n")
		if m.uvStatus.Version != "" {
			content.WriteString(fmt.Sprintf("   Version: %s\n", m.uvStatus.Version))
		}
		if m.uvStatus.Path != "" {
			content.WriteString(fmt.Sprintf("   Path: %s\n", m.uvStatus.Path))
		}
	} else {
		content.WriteString("❌ UV is not installed\n")
		content.WriteString("   Press 'i' to install UV\n\n")

		// Show install command
		cmd, err := m.uvInstaller.GetInstallCommand()
		if err == nil {
			content.WriteString("Install command:\n")
			content.WriteString(fmt.Sprintf("   %s\n", cmd))
		}
	}

	// Show recent messages
	if len(m.messages) > 0 {
		content.WriteString("\nRecent Messages:\n")
		for _, msg := range m.messages {
			content.WriteString(fmt.Sprintf("• %s\n", msg))
		}
	}

	return content.String()
}

// renderPythonPanel renders the Python management panel
func (m Model) renderPythonPanel() string {
	content := "Python Version Management\n\n"
	if !m.uvStatus.Installed {
		content += "UV must be installed first to manage Python versions."
	} else {
		content += "Available Python versions:\n"
		content += "• Python 3.12 (not implemented yet)\n"
		content += "• Python 3.11 (not implemented yet)\n"
		content += "• Python 3.10 (not implemented yet)\n"
	}
	return content
}

// renderProjectPanel renders the project management panel
func (m Model) renderProjectPanel() string {
	content := "Project Management\n\n"
	if !m.uvStatus.Installed {
		content += "UV must be installed first to manage projects."
	} else {
		content += "Current project status:\n"
		content += "• No project detected in current directory\n"
		content += "• Use 'uv init' to create a new project (not implemented yet)\n"
	}
	return content
}

// renderEnvironmentPanel renders the environment management panel
func (m Model) renderEnvironmentPanel() string {
	content := "Environment Management\n\n"
	if !m.uvStatus.Installed {
		content += "UV must be installed first to manage environments."
	} else {
		content += "Virtual environments:\n"
		content += "• No virtual environments found\n"
		content += "• Use 'uv venv' to create environments (not implemented yet)\n"
	}
	return content
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	osInfo := fmt.Sprintf("OS: %s", getOS())
	panelInfo := fmt.Sprintf("Panel: %d/%d", int(m.activePanel)+1, len(m.panels))

	return statusStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			osInfo,
			strings.Repeat(" ", max(0, m.width-len(osInfo)-len(panelInfo)-4)),
			panelInfo,
		),
	)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	program := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

package app

import (
	"fmt"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"uvui/internal/types"
	"uvui/internal/ui"
	"uvui/internal/ui/panels"
)

// Init initializes the application
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		CheckUVStatus(m.UVInstaller),
		tea.EnterAltScreen,
	)
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.State.Width = msg.Width
		m.State.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case ui.UVInstalledMsg:
		return m.handleUVInstalledMsg(msg)

	case ui.PythonVersionsLoadedMsg:
		return m.handlePythonVersionsLoadedMsg(msg)

	case ui.PythonOperationMsg:
		return m.handlePythonOperationMsg(msg)
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		return m.handleTabNavigation(1)

	case "shift+tab":
		return m.handleTabNavigation(-1)

	case "up", "k":
		return m.handleVerticalNavigation(-1)

	case "down", "j":
		return m.handleVerticalNavigation(1)

	case "enter":
		return m.handleEnterKey()

	case "d", "delete":
		return m.handleDeleteKey()

	case "p":
		return m.handlePinKey()

	case "i":
		return m.handleInstallRefresh()

	case "r":
		return m.handleRefresh()

	case "?":
		return m.handleHelp()
	}

	return m, nil
}

// handleTabNavigation handles tab and shift+tab navigation
func (m *Model) handleTabNavigation(direction int) (tea.Model, tea.Cmd) {
	if direction > 0 {
		m.State.ActivePanel = types.Panel((int(m.State.ActivePanel) + 1) % len(m.State.Panels))
	} else {
		if m.State.ActivePanel == 0 {
			m.State.ActivePanel = types.Panel(len(m.State.Panels) - 1)
		} else {
			m.State.ActivePanel = types.Panel(int(m.State.ActivePanel) - 1)
		}
	}

	// Load Python versions when entering Python panel
	if m.State.ActivePanel == types.PythonPanel && m.State.UVStatus.Installed && !m.State.PythonVersions.Loading {
		m.State.PythonVersions.Loading = true
		return m, LoadPythonVersions(m.PythonManager)
	}

	return m, nil
}

// handleVerticalNavigation handles up/down arrow navigation within panels
func (m *Model) handleVerticalNavigation(direction int) (tea.Model, tea.Cmd) {
	if m.State.ActivePanel == types.PythonPanel {
		// Merge available and installed for display order
		allVersions := panels.MergePythonVersions(m.State.PythonVersions.Available, m.State.PythonVersions.Installed)
		if len(allVersions) > 0 {
			if direction < 0 && m.State.PythonVersions.Selected > 0 {
				m.State.PythonVersions.Selected--
			} else if direction > 0 && m.State.PythonVersions.Selected < len(allVersions)-1 {
				m.State.PythonVersions.Selected++
			}
		}
	}
	return m, nil
}

// handleEnterKey handles enter key press
func (m *Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.State.ActivePanel == types.PythonPanel && m.State.UVStatus.Installed && !m.State.Operation.InProgress {
		if selectedVersion := m.GetSelectedPythonVersion(); selectedVersion != nil && !selectedVersion.Installed {
			m.SetOperation("install", selectedVersion.Version, true)
			m.AddMessage(fmt.Sprintf("Installing Python %s...", selectedVersion.Version))
			return m, InstallPythonVersion(m.PythonManager, selectedVersion.Version)
		}
	}
	return m, nil
}

// handleDeleteKey handles delete key press
func (m *Model) handleDeleteKey() (tea.Model, tea.Cmd) {
	if m.State.ActivePanel == types.PythonPanel && m.State.UVStatus.Installed && !m.State.Operation.InProgress {
		if selectedVersion := m.GetSelectedPythonVersion(); selectedVersion != nil && selectedVersion.Installed && !selectedVersion.Current {
			m.SetOperation("uninstall", selectedVersion.Version, true)
			m.AddMessage(fmt.Sprintf("Uninstalling Python %s...", selectedVersion.Version))
			return m, UninstallPythonVersion(m.PythonManager, selectedVersion.Version)
		}
	}
	return m, nil
}

// handlePinKey handles pin key press
func (m *Model) handlePinKey() (tea.Model, tea.Cmd) {
	if m.State.ActivePanel == types.PythonPanel && m.State.UVStatus.Installed && !m.State.Operation.InProgress {
		if selectedVersion := m.GetSelectedPythonVersion(); selectedVersion != nil && selectedVersion.Installed {
			m.SetOperation("pin", selectedVersion.Version, true)
			m.AddMessage(fmt.Sprintf("Pinning Python %s...", selectedVersion.Version))
			return m, PinPythonVersion(m.PythonManager, selectedVersion.Version)
		}
	}
	return m, nil
}

// handleInstallRefresh handles install/refresh key press
func (m *Model) handleInstallRefresh() (tea.Model, tea.Cmd) {
	switch m.State.ActivePanel {
	case types.StatusPanel:
		if !m.State.UVStatus.Installed && !m.State.Installing {
			m.State.Installing = true
			m.AddMessage("Installing UV...")
			return m, InstallUV(m.UVInstaller)
		}
	case types.PythonPanel:
		if m.State.UVStatus.Installed && !m.State.PythonVersions.Loading {
			m.State.PythonVersions.Loading = true
			m.AddMessage("Refreshing Python versions...")
			return m, LoadPythonVersions(m.PythonManager)
		}
	}
	return m, nil
}

// handleRefresh handles refresh key press
func (m *Model) handleRefresh() (tea.Model, tea.Cmd) {
	switch m.State.ActivePanel {
	case types.StatusPanel:
		m.AddMessage("Refreshing UV status...")
		return m, CheckUVStatus(m.UVInstaller)
	case types.PythonPanel:
		if m.State.UVStatus.Installed && !m.State.PythonVersions.Loading {
			m.State.PythonVersions.Loading = true
			m.AddMessage("Refreshing Python versions...")
			return m, LoadPythonVersions(m.PythonManager)
		}
	}
	return m, nil
}

// handleHelp handles help key press
func (m *Model) handleHelp() (tea.Model, tea.Cmd) {
	helpText := m.getCurrentPanelHelp()
	m.AddMessage(fmt.Sprintf("Help: %s", helpText))
	return m, nil
}

// getCurrentPanelHelp returns help text for the current panel
func (m *Model) getCurrentPanelHelp() string {
	switch m.State.ActivePanel {
	case types.StatusPanel:
		return panels.GetStatusPanelHelp()
	case types.PythonPanel:
		return panels.GetPythonPanelHelp()
	case types.ProjectPanel:
		return panels.GetProjectPanelHelp()
	case types.EnvironmentPanel:
		return panels.GetEnvironmentPanelHelp()
	default:
		return "q=quit, tab=next panel, shift+tab=prev panel"
	}
}

// Message handlers
func (m *Model) handleUVInstalledMsg(msg ui.UVInstalledMsg) (tea.Model, tea.Cmd) {
	m.State.Installing = false
	if msg.Success {
		m.State.UVStatus.Installed = true
		m.State.UVStatus.Version = msg.Version
		m.AddMessage("UV installation completed successfully!")
	} else if msg.Error != nil {
		m.AddMessage(fmt.Sprintf("UV installation failed: %v", msg.Error))
	} else {
		m.State.UVStatus.Installed = false
		m.AddMessage("UV is not installed")
	}
	return m, nil
}

func (m *Model) handlePythonVersionsLoadedMsg(msg ui.PythonVersionsLoadedMsg) (tea.Model, tea.Cmd) {
	m.State.PythonVersions.Loading = false
	if msg.Error != nil {
		m.AddMessage(fmt.Sprintf("Failed to load Python versions: %v", msg.Error))
		return m, nil
	}

	m.UpdatePythonVersions(msg.Available, msg.Installed)
	m.AddMessage(fmt.Sprintf("Loaded %d Python versions", len(msg.Available)))
	return m, nil
}

func (m *Model) handlePythonOperationMsg(msg ui.PythonOperationMsg) (tea.Model, tea.Cmd) {
	m.CompleteOperation(msg.Success, msg.Error)

	if msg.Success {
		m.AddMessage(fmt.Sprintf("Successfully %sed Python %s", msg.Operation, msg.Target))
		// Reload Python versions after successful operation
		m.State.PythonVersions.Loading = true
		return m, LoadPythonVersions(m.PythonManager)
	} else {
		m.AddMessage(fmt.Sprintf("Failed to %s Python %s: %v", msg.Operation, msg.Target, msg.Error))
	}

	return m, nil
}

// View renders the application UI
func (m *Model) View() string {
	if m.State.Width == 0 {
		return "Initializing..."
	}

	// Header
	header := ui.TitleStyle.Render("UV Package Manager CLI")

	// Tabs
	tabs := m.renderTabs()

	// Main content based on active panel
	content := m.renderActivePanel()

	// Status bar
	statusBar := m.renderStatusBar()

	// Help line
	helpLine := ui.HelpStyle.Render("? for help | q to quit | tab to navigate")

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
func (m *Model) renderTabs() string {
	tabNames := []string{"Status", "Python", "Project", "Environment"}
	var tabs []string

	for i, name := range tabNames {
		if types.Panel(i) == m.State.ActivePanel {
			tabs = append(tabs, ui.ActiveTabStyle.Render(name))
		} else {
			tabs = append(tabs, ui.InactiveTabStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderActivePanel renders the content of the active panel
func (m *Model) renderActivePanel() string {
	style := ui.PanelStyle
	if m.State.ActivePanel == types.StatusPanel {
		style = ui.ActivePanelStyle
	}

	var content string
	switch m.State.ActivePanel {
	case types.StatusPanel:
		content = panels.RenderStatusPanel(m.State, m.UVInstaller)
	case types.PythonPanel:
		style = ui.ActivePanelStyle
		content = panels.RenderPythonPanel(m.State)
	case types.ProjectPanel:
		content = panels.RenderProjectPanel(m.State)
	case types.EnvironmentPanel:
		content = panels.RenderEnvironmentPanel(m.State)
	default:
		content = "Unknown panel"
	}

	return style.Render(content)
}

// renderStatusBar renders the bottom status bar
func (m *Model) renderStatusBar() string {
	osInfo := fmt.Sprintf("OS: %s", getOS())
	panelInfo := fmt.Sprintf("Panel: %d/%d", int(m.State.ActivePanel)+1, len(m.State.Panels))

	return ui.StatusStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			osInfo,
			strings.Repeat(" ", max(0, m.State.Width-len(osInfo)-len(panelInfo)-4)),
			panelInfo,
		),
	)
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	case "freebsd":
		return "FreeBSD"
	case "openbsd":
		return "OpenBSD"
	case "netbsd":
		return "NetBSD"
	default:
		return runtime.GOOS
	}
}

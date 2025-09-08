// Package app provides the core application logic.
package app

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"uvui/internal/services"
	"uvui/internal/ui"
)

// handleInitConfig creates a new keybindings.json file with default values.
func (m *Model) handleInitConfig() (tea.Model, tea.Cmd) {
	_, err := os.Stat("keybindings.json")
	if err == nil {
		m.AddMessage("keybindings.json already exists.")
		return m, nil
	}

	if os.IsNotExist(err) {
		err := InitConfig("keybindings.json")
		if err != nil {
			m.AddMessage("Error creating keybindings.json: " + err.Error())
			return m, nil
		}
		m.AddMessage("keybindings.json created successfully.")
		return m, nil
	}

	m.AddMessage("Error checking keybindings.json: " + err.Error())
	return m, nil
}

// CheckUVStatus checks the UV installation status.
func CheckUVStatus(installer services.UVInstallerInterface) tea.Cmd {
	return func() tea.Msg {
		installed, version, err := installer.IsInstalled()
		if err != nil {
			return ui.UVInstalledMsg{Success: false, Error: err}
		}

		return ui.UVInstalledMsg{
			Success: installed,
			Version: version,
			Error:   nil,
		}
	}
}

// InstallUV installs UV.
func InstallUV(installer services.UVInstallerInterface) tea.Cmd {
	return func() tea.Msg {
		err := installer.Install()
		return ui.UVInstalledMsg{Success: err == nil, Error: err}
	}
}

// LoadPythonVersions loads Python version information.
func LoadPythonVersions(manager services.PythonManagerInterface) tea.Cmd {
	return func() tea.Msg {
		available, err := manager.ListAvailable()
		if err != nil {
			return ui.PythonVersionsLoadedMsg{Error: err}
		}

		installed, err := manager.ListInstalled()
		if err != nil {
			return ui.PythonVersionsLoadedMsg{Available: available, Error: err}
		}

		return ui.PythonVersionsLoadedMsg{
			Available: available,
			Installed: installed,
		}
	}
}

// InstallPythonVersion installs a Python version.
func InstallPythonVersion(manager services.PythonManagerInterface, version string) tea.Cmd {
	return func() tea.Msg {
		err := manager.Install(version)
		return ui.PythonOperationMsg{
			Operation: "install",
			Target:    version,
			Success:   err == nil,
			Error:     err,
		}
	}
}

// UninstallPythonVersion uninstalls a Python version.
func UninstallPythonVersion(manager services.PythonManagerInterface, version string) tea.Cmd {
	return func() tea.Msg {
		err := manager.Uninstall(version)
		return ui.PythonOperationMsg{
			Operation: "uninstall",
			Target:    version,
			Success:   err == nil,
			Error:     err,
		}
	}
}

// PinPythonVersion pins a Python version.
func PinPythonVersion(manager services.PythonManagerInterface, version string) tea.Cmd {
	return func() tea.Msg {
		err := manager.Pin(version)
		return ui.PythonOperationMsg{
			Operation: "pin",
			Target:    version,
			Success:   err == nil,
			Error:     err,
		}
	}
}

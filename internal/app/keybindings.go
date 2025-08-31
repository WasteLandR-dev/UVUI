// Package app provides the core application logic.
package app

import (
	"encoding/json"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Keybindings holds the keybindings for the application
type Keybindings struct {
	Quit           []string `json:"quit"`
	NextPanel      []string `json:"next_panel"`
	PrevPanel      []string `json:"prev_panel"`
	NavUp          []string `json:"navigation_up"`
	NavDown        []string `json:"navigation_down"`
	Install        []string `json:"install"`
	Delete         []string `json:"delete"`
	Pin            []string `json:"pin"`
	Refresh        []string `json:"refresh"`
	Help           []string `json:"help"`
	Sync           []string `json:"sync"`
	Lock           []string `json:"lock"`
	ToggleView     []string `json:"toggle_view"`
	InitApp        []string `json:"init_app"`
	InitNew        []string `json:"init_new"`
	InstallRefresh []string `json:"install_refresh"`
}

// Config holds the application configuration
type Config struct {
	Keybindings Keybindings `json:"keybindings"`
}

// LoadConfig loads the configuration from the given path
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case contains(m.Config.Keybindings.Quit, msg.String()):
		return m, tea.Quit
	case contains(m.Config.Keybindings.NextPanel, msg.String()):
		return m.handleTabNavigation(1)
	case contains(m.Config.Keybindings.PrevPanel, msg.String()):
		return m.handleTabNavigation(-1)
	case contains(m.Config.Keybindings.NavUp, msg.String()):
		return m.handleVerticalNavigation(-1)
	case contains(m.Config.Keybindings.NavDown, msg.String()):
		return m.handleVerticalNavigation(1)
	case contains(m.Config.Keybindings.Install, msg.String()):
		return m.handleEnterKey()
	case contains(m.Config.Keybindings.Delete, msg.String()):
		return m.handleDeleteKey()
	case contains(m.Config.Keybindings.Pin, msg.String()):
		return m.handlePinKey()
	case contains(m.Config.Keybindings.Refresh, msg.String()):
		return m.handleRefresh()
	case contains(m.Config.Keybindings.Help, msg.String()):
		return m.handleHelp()
	case contains(m.Config.Keybindings.Sync, msg.String()):
		return m.handleSyncKey()
	case contains(m.Config.Keybindings.Lock, msg.String()):
		return m.handleLockOrLibKey()
	case contains(m.Config.Keybindings.ToggleView, msg.String()):
		return m.handleToggleKey()
	case contains(m.Config.Keybindings.InitApp, msg.String()):
		return m.handleAppKey()
	case contains(m.Config.Keybindings.InitNew, msg.String()):
		return m.handleNewProjectKey()
	case contains(m.Config.Keybindings.InstallRefresh, msg.String()):
		return m.handleInstallRefresh()
	}

	return m, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

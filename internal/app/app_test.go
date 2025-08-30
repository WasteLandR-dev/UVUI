package app

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"uvui/internal/types"
	"uvui/internal/ui"
)

// Simple test helpers

func TestNewModel(t *testing.T) {
	model := NewModel()

	assert.NotNil(t, model)
	assert.NotNil(t, model.State)
	assert.NotNil(t, model.UVInstaller)
	assert.NotNil(t, model.PythonManager)
	assert.NotNil(t, model.CommandExecutor)

	// Check initial state
	assert.Equal(t, types.StatusPanel, model.State.ActivePanel)
	assert.Len(t, model.State.Panels, 4)
	assert.Equal(t, types.StatusPanel, model.State.Panels[0])
	assert.Equal(t, types.PythonPanel, model.State.Panels[1])
	assert.Equal(t, types.ProjectPanel, model.State.Panels[2])
	assert.Equal(t, types.EnvironmentPanel, model.State.Panels[3])
}

func TestModel_AddMessage(t *testing.T) {
	model := NewModel()

	// Add messages
	model.AddMessage("Test message 1")
	model.AddMessage("Test message 2")

	assert.Len(t, model.State.Messages, 2)
	assert.Equal(t, "Test message 1", model.State.Messages[0])
	assert.Equal(t, "Test message 2", model.State.Messages[1])

	// Test message limit (should keep only last 10)
	for i := 0; i < 15; i++ {
		model.AddMessage(fmt.Sprintf("Message %d", i))
	}

	assert.Len(t, model.State.Messages, 10)
	assert.Equal(t, "Message 5", model.State.Messages[0])
	assert.Equal(t, "Message 14", model.State.Messages[9])
}

func TestModel_GetSelectedPythonVersion(t *testing.T) {
	model := NewModel()

	// Test with no versions
	selected := model.GetSelectedPythonVersion()
	assert.Nil(t, selected)

	// Test with versions
	model.State.PythonVersions.Available = []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	model.State.PythonVersions.Installed = []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true},
	}

	selected = model.GetSelectedPythonVersion()
	assert.NotNil(t, selected)
	assert.Equal(t, "3.12.0", selected.Version)
}

func TestModel_UpdatePythonVersions(t *testing.T) {
	model := NewModel()

	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	installed := []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true},
	}

	model.UpdatePythonVersions(available, installed)

	assert.Equal(t, available, model.State.PythonVersions.Available)
	assert.Equal(t, installed, model.State.PythonVersions.Installed)
	assert.False(t, model.State.PythonVersions.Loading)
}

func TestModel_SetOperation(t *testing.T) {
	model := NewModel()

	model.SetOperation("install", "3.12.0", true)

	assert.True(t, model.State.Operation.InProgress)
	assert.Equal(t, "install", model.State.Operation.Operation)
	assert.Equal(t, "3.12.0", model.State.Operation.Target)
}

func TestModel_CompleteOperation(t *testing.T) {
	model := NewModel()

	model.SetOperation("install", "3.12.0", true)
	model.CompleteOperation(true, nil)

	assert.False(t, model.State.Operation.InProgress)
	assert.True(t, model.State.Operation.Success)
	assert.Nil(t, model.State.Operation.Error)
}

func TestModel_Init(t *testing.T) {
	model := NewModel()

	cmd := model.Init()

	assert.NotNil(t, cmd)
}

func TestModel_Update_WindowSize(t *testing.T) {
	model := NewModel()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, cmd := model.Update(msg)

	assert.Equal(t, 100, model.State.Width)
	assert.Equal(t, 50, model.State.Height)
	assert.Nil(t, cmd)

	// Ensure we return the same model
	assert.Equal(t, model, updatedModel)
}

func TestModel_Update_KeyPress(t *testing.T) {
	model := NewModel()

	// Test quit key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(msg)

	assert.NotNil(t, cmd)

	// Ensure we return the same model
	assert.Equal(t, model, updatedModel)
}

func TestModel_handleTabNavigation(t *testing.T) {
	model := NewModel()

	// Test forward navigation
	updatedModel, cmd := model.handleTabNavigation(1)

	assert.Equal(t, types.PythonPanel, model.State.ActivePanel)
	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)

	// Test backward navigation
	updatedModel, cmd = model.handleTabNavigation(-1)

	assert.Equal(t, types.StatusPanel, model.State.ActivePanel)
	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)
}

func TestModel_handleVerticalNavigation(t *testing.T) {
	model := NewModel()
	model.State.ActivePanel = types.PythonPanel

	// Add some Python versions
	model.State.PythonVersions.Available = []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
		{Version: "3.10.0", Installed: false},
	}
	model.State.PythonVersions.Installed = []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true},
	}

	// Test down navigation
	updatedModel, cmd := model.handleVerticalNavigation(1)

	assert.Equal(t, 1, model.State.PythonVersions.Selected)
	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)

	// Test up navigation
	updatedModel, cmd = model.handleVerticalNavigation(-1)

	assert.Equal(t, 0, model.State.PythonVersions.Selected)
	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)
}

func TestModel_handleEnterKey(t *testing.T) {
	model := NewModel()
	model.State.ActivePanel = types.PythonPanel
	model.State.Installed = true

	// Add Python versions
	model.State.PythonVersions.Available = []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
	}
	model.State.PythonVersions.Installed = []types.PythonVersion{}

	updatedModel, cmd := model.handleEnterKey()

	assert.NotNil(t, cmd)
	assert.Equal(t, model, updatedModel)
}

func TestModel_handleHelp(t *testing.T) {
	model := NewModel()

	updatedModel, cmd := model.handleHelp()

	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)

	// Check that help message was added
	assert.Len(t, model.State.Messages, 1)
	assert.Contains(t, model.State.Messages[0], "Help:")
}

func TestModel_View(t *testing.T) {
	model := NewModel()
	model.State.Width = 100
	model.State.Height = 50

	view := model.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "UV Package Manager CLI")
	assert.Contains(t, view, "Status")
	assert.Contains(t, view, "Python")
	assert.Contains(t, view, "Project")
	assert.Contains(t, view, "Environment")
}

func TestModel_renderTabs(t *testing.T) {
	model := NewModel()

	tabs := model.renderTabs()

	assert.NotEmpty(t, tabs)
	assert.Contains(t, tabs, "Status")
	assert.Contains(t, tabs, "Python")
	assert.Contains(t, tabs, "Project")
	assert.Contains(t, tabs, "Environment")
}

func TestModel_renderStatusBar(t *testing.T) {
	model := NewModel()
	model.State.Width = 100
	model.State.Height = 50

	statusBar := model.renderStatusBar()

	assert.NotEmpty(t, statusBar)
	assert.Contains(t, statusBar, "OS:")
	assert.Contains(t, statusBar, "Panel:")
}

func TestModel_getCurrentPanelHelp(t *testing.T) {
	model := NewModel()

	// Test Status panel help
	help := model.getCurrentPanelHelp()
	assert.NotEmpty(t, help)

	// Test Python panel help
	model.State.ActivePanel = types.PythonPanel
	help = model.getCurrentPanelHelp()
	assert.NotEmpty(t, help)

	// Test Project panel help
	model.State.ActivePanel = types.ProjectPanel
	help = model.getCurrentPanelHelp()
	assert.NotEmpty(t, help)

	// Test Environment panel help
	model.State.ActivePanel = types.EnvironmentPanel
	help = model.getCurrentPanelHelp()
	assert.NotEmpty(t, help)
}

func TestModel_handleUVInstalledMsg(t *testing.T) {
	model := NewModel()

	msg := ui.UVInstalledMsg{Success: true, Version: "1.0.0"}
	updatedModel, cmd := model.handleUVInstalledMsg(msg)

	assert.True(t, model.State.Installed)
	assert.Equal(t, "1.0.0", model.State.Version)
	assert.False(t, model.State.Installing)
	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)

	// Check message was added
	assert.Len(t, model.State.Messages, 1)
	assert.Contains(t, model.State.Messages[0], "UV installation completed successfully!")
}

func TestModel_handlePythonVersionsLoadedMsg(t *testing.T) {
	model := NewModel()

	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
	}
	installed := []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true},
	}

	msg := ui.PythonVersionsLoadedMsg{
		Available: available,
		Installed: installed,
	}

	updatedModel, cmd := model.handlePythonVersionsLoadedMsg(msg)

	assert.False(t, model.State.PythonVersions.Loading)
	assert.Equal(t, available, model.State.PythonVersions.Available)
	assert.Equal(t, installed, model.State.PythonVersions.Installed)
	assert.Nil(t, cmd)
	assert.Equal(t, model, updatedModel)

	// Check message was added
	assert.Len(t, model.State.Messages, 1)
	assert.Contains(t, model.State.Messages[0], "Loaded 1 Python versions")
}

func TestModel_handlePythonOperationMsg(t *testing.T) {
	model := NewModel()

	msg := ui.PythonOperationMsg{
		Operation: "install",
		Target:    "3.12.0",
		Success:   true,
		Error:     nil,
	}

	updatedModel, cmd := model.handlePythonOperationMsg(msg)

	assert.False(t, model.State.Operation.InProgress)
	assert.NotNil(t, cmd)
	assert.Equal(t, model, updatedModel)

	// Check message was added
	assert.Len(t, model.State.Messages, 1)
	assert.Contains(t, model.State.Messages[0], "Successfully installed Python 3.12.0")
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 5, maxInt(5, 3))
	assert.Equal(t, 5, maxInt(3, 5))
	assert.Equal(t, 5, maxInt(5, 5))
}

func TestGetOS(t *testing.T) {
	os := getOS()
	assert.NotEmpty(t, os)
}

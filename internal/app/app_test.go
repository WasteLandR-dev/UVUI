package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"uvui/internal/services"
	"uvui/internal/types"
	"uvui/internal/ui"
)

// mockCommandExecutor is a mock implementation of the CommandExecutorInterface.
type mockCommandExecutor struct {
	services.CommandExecutorInterface
	IsUVAvailableFunc func() bool
	ExecuteFunc       func(command string, args ...string) ([]byte, error)
}

func (m *mockCommandExecutor) IsUVAvailable() bool {
	if m.IsUVAvailableFunc != nil {
		return m.IsUVAvailableFunc()
	}
	return true
}

func (m *mockCommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command, args...)
	}
	return nil, nil
}

// newTestModel creates a new model with mock services for testing.
func newTestModel() *Model {
	executor := &mockCommandExecutor{}
	uvInstaller := services.NewUVInstaller(executor)
	pythonManager := services.NewPythonManager(executor)
	projectManager := services.NewProjectManager(executor)
	return NewModel(uvInstaller, pythonManager, projectManager, executor)
}

func TestGetDirection(t *testing.T) {
	assert.Equal(t, -1, getDirection("up"))
	assert.Equal(t, 1, getDirection("down"))
	assert.Equal(t, 0, getDirection("other"))
}

func TestInit(t *testing.T) {
	m := newTestModel()
	cmd := m.Init()
	assert.NotNil(t, cmd)
}

func TestUpdate(t *testing.T) {
	m := newTestModel()
	model, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	assert.NotNil(t, model)
	assert.Nil(t, cmd)

	// Test key press
	model, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleTextInput(t *testing.T) {
	m := newTestModel()
	m.InputMode = InputModeProjectName
	m.TextInput.SetValue("test-project")

	model, cmd := m.handleTextInput(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
	assert.Equal(t, InputModeNone, m.InputMode)
}

func TestHandleTabNavigation(t *testing.T) {
	m := newTestModel()
	initialPanel := m.State.ActivePanel

	model, cmd := m.handleTabNavigation(1)
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
	assert.NotEqual(t, initialPanel, m.State.ActivePanel)
}

func TestHandleVerticalNavigation(t *testing.T) {
	m := newTestModel()
	m.State.PythonVersions.Available = []types.PythonVersion{{Version: "3.12.1"}}
	m.State.PythonVersions.Selected = 0

	model, cmd := m.handleVerticalNavigation(1)
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, m.State.PythonVersions.Selected)
}

func TestHandleEnterKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.PythonPanel
	m.State.Installed = true
	m.State.PythonVersions.Available = []types.PythonVersion{{Version: "3.12.1", Installed: false}}

	model, cmd := m.handleEnterKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleDeleteKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.PythonPanel
	m.State.Installed = true
	m.State.PythonVersions.Installed = []types.PythonVersion{{Version: "3.12.1", Installed: true}}

	model, cmd := m.handleDeleteKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandlePinKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.PythonPanel
	m.State.Installed = true
	m.State.PythonVersions.Installed = []types.PythonVersion{{Version: "3.12.1", Installed: true}}

	model, cmd := m.handlePinKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleInstallRefresh(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.StatusPanel

	model, cmd := m.handleInstallRefresh()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleRefresh(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.StatusPanel

	model, cmd := m.handleRefresh()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleHelp(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handleHelp()
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
	assert.NotEmpty(t, m.State.Messages)
}

func TestGetCurrentPanelHelp(t *testing.T) {
	m := newTestModel()
	helpText := m.getCurrentPanelHelp()
	assert.NotEmpty(t, helpText)
}

func TestHandleUVInstalledMsg(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handleUVInstalledMsg(ui.UVInstalledMsg{Success: true, Version: "0.1.0"})
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
	assert.True(t, m.State.Installed)
}

func TestHandlePythonVersionsLoadedMsg(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handlePythonVersionsLoadedMsg(ui.PythonVersionsLoadedMsg{Available: []types.PythonVersion{{Version: "3.12.1"}}})
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
	assert.NotEmpty(t, m.State.PythonVersions.Available)
}

func TestHandlePythonOperationMsg(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handlePythonOperationMsg(ui.PythonOperationMsg{Success: true, Operation: "install", Target: "3.12.1"})
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestView(t *testing.T) {
	m := newTestModel()
	m.State.Width = 80
	view := m.View()
	assert.NotEmpty(t, view)
}

func TestRenderTabs(t *testing.T) {
	m := newTestModel()
	tabs := m.renderTabs()
	assert.NotEmpty(t, tabs)
}

func TestRenderActivePanel(t *testing.T) {
	m := newTestModel()
	panel := m.renderActivePanel()
	assert.NotEmpty(t, panel)
}

func TestRenderStatusBar(t *testing.T) {
	m := newTestModel()
	m.State.Width = 80
	statusBar := m.renderStatusBar()
	assert.NotEmpty(t, statusBar)
}

func TestHandleSyncKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.ProjectPanel
	m.State.Installed = true
	m.State.ProjectState.Status = &types.ProjectStatus{IsProject: true}

	model, cmd := m.handleSyncKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleLockOrLibKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.ProjectPanel
	m.State.Installed = true
	m.State.ProjectState.Status = &types.ProjectStatus{IsProject: true}

	model, cmd := m.handleLockOrLibKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleToggleKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.ProjectPanel
	m.State.ProjectState.Status = &types.ProjectStatus{IsProject: true}

	model, cmd := m.handleToggleKey()
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
}

func TestHandleAppKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.ProjectPanel
	m.State.Installed = true

	model, cmd := m.handleAppKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleNewProjectKey(t *testing.T) {
	m := newTestModel()
	m.State.ActivePanel = types.ProjectPanel
	m.State.Installed = true

	model, cmd := m.handleNewProjectKey()
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleProjectStatusLoadedMsg(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handleProjectStatusLoadedMsg(ui.ProjectStatusLoadedMsg{Status: &types.ProjectStatus{IsProject: true}})
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestHandleProjectDependenciesLoadedMsg(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handleProjectDependenciesLoadedMsg(ui.ProjectDependenciesLoadedMsg{Dependencies: []types.ProjectDependency{{Name: "requests"}}})
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
}

func TestHandleProjectOperationMsg(t *testing.T) {
	m := newTestModel()
	model, cmd := m.handleProjectOperationMsg(ui.ProjectOperationMsg{Success: true, Operation: "init"})
	assert.NotNil(t, model)
	assert.NotNil(t, cmd)
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 5, maxInt(3, 5))
	assert.Equal(t, 5, maxInt(5, 3))
}

func TestGetOS(t *testing.T) {
	assert.NotEmpty(t, getOS())
}

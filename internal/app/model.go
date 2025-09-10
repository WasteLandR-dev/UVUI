// Package app provides the core application logic.
package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	"uvui/internal/services"
	"uvui/internal/types"
	"uvui/internal/ui/panels"
)

// Model represents the application state and dependencies.
type Model struct {
	State           *panels.AppState
	Config          *Config
	UVInstaller     services.UVInstallerInterface
	PythonManager   services.PythonManagerInterface
	ProjectManager  services.ProjectManagerInterface
	CommandExecutor services.CommandExecutorInterface
	TextInput       textinput.Model
	InputMode       InputMode
}

// NewModel creates a new application model.
func NewModel(uvInstaller services.UVInstallerInterface, pythonManager services.PythonManagerInterface, projectManager services.ProjectManagerInterface, commandExecutor services.CommandExecutorInterface) *Model {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	ti := textinput.New()
	ti.Placeholder = "Project Name"
	ti.Focus()

	state := &panels.AppState{
		ActivePanel: types.StatusPanel,
		Panels: []types.Panel{
			types.StatusPanel,
			types.PythonPanel,
			types.ProjectPanel,
			types.EnvironmentPanel,
		},
		PythonVersions: panels.PythonVersions{
			Available: []types.PythonVersion{},
			Installed: []types.PythonVersion{},
			Selected:  0,
			Loading:   false,
		},
		ProjectState: panels.ProjectState{
			Status:         nil,
			Dependencies:   []types.ProjectDependency{},
			DependencyTree: nil,
			Selected:       0,
			Loading:        false,
			ShowTree:       false,
		},
		Messages: []string{},
	}

	m := &Model{
		State:           state,
		Config:          config,
		UVInstaller:     uvInstaller,
		PythonManager:   pythonManager,
		ProjectManager:  projectManager,
		CommandExecutor: commandExecutor,
		TextInput:       ti,
		InputMode:       InputModeNone,
	}

	if config.KeybindingsNotFound {
		m.AddMessage("keybindings.json not found. Using default keybindings.")
		m.AddMessage("Press 'c' to create a default keybindings.json file.")
		m.AddMessage("Create keybindings.json in the same directory as the executable to customize.")
	}

	return m
}

// AddMessage adds a message to the message list.
func (m *Model) AddMessage(msg string) {
	m.State.Messages = append(m.State.Messages, msg)
	// Keep only last 10 messages
	if len(m.State.Messages) > 10 {
		m.State.Messages = m.State.Messages[1:]
	}
}

// GetSelectedPythonVersion returns the currently selected Python version.
func (m *Model) GetSelectedPythonVersion() *types.PythonVersion {
	// Get the merged list to match what's displayed
	merged := m.GetMergedPythonVersions()
	if len(merged) == 0 {
		return nil
	}

	// Don't modify selection here - just return nil if out of bounds
	if m.State.PythonVersions.Selected >= len(merged) || m.State.PythonVersions.Selected < 0 {
		return nil
	}

	return &merged[m.State.PythonVersions.Selected]
}

// ValidateAndFixSelection ensures the selection index is within bounds.
func (m *Model) ValidateAndFixSelection() {
	merged := m.GetMergedPythonVersions()
	if len(merged) == 0 {
		m.State.PythonVersions.Selected = 0
		return
	}

	if m.State.PythonVersions.Selected >= len(merged) {
		m.State.PythonVersions.Selected = 0
	}
	if m.State.PythonVersions.Selected < 0 {
		m.State.PythonVersions.Selected = 0
	}
}

// GetMergedPythonVersions returns the merged list of available and installed versions.
func (m *Model) GetMergedPythonVersions() []types.PythonVersion {
	// Import the panels package to use MergePythonVersions
	// This is a bit of a circular dependency, but necessary for consistency
	return panels.MergePythonVersions(m.State.PythonVersions.Available, m.State.PythonVersions.Installed)
}

// GetMergedPythonVersionsCount returns the count of merged Python versions.
func (m *Model) GetMergedPythonVersionsCount() int {
	return len(m.GetMergedPythonVersions())
}

// UpdatePythonVersions updates the Python versions in the state.
func (m *Model) UpdatePythonVersions(available, installed []types.PythonVersion) {
	m.State.PythonVersions.Available = available
	m.State.PythonVersions.Installed = installed
	m.State.PythonVersions.Loading = false

	// Validate and fix selection after updating versions
	m.ValidateAndFixSelection()
}

// UpdateProjectStatus updates the project status in the state.
func (m *Model) UpdateProjectStatus(status *types.ProjectStatus) {
	m.State.ProjectState.Status = status
	m.State.ProjectState.Loading = false
}

// UpdateProjectDependencies updates project dependencies and tree.
func (m *Model) UpdateProjectDependencies(deps []types.ProjectDependency, tree *types.DependencyTree) {
	m.State.ProjectState.Dependencies = deps
	m.State.ProjectState.DependencyTree = tree
}

// ToggleTreeView toggles between dependency list and tree view.
func (m *Model) ToggleTreeView() {
	m.State.ProjectState.ShowTree = !m.State.ProjectState.ShowTree
}

// SetProjectLoading sets the project loading state.
func (m *Model) SetProjectLoading(loading bool) {
	m.State.ProjectState.Loading = loading
}

// SetOperation sets the current operation status.
func (m *Model) SetOperation(operation, target string, inProgress bool) {
	m.State.Operation = types.OperationStatus{
		InProgress: inProgress,
		Operation:  operation,
		Target:     target,
	}
}

// CompleteOperation completes the current operation.
func (m *Model) CompleteOperation(success bool, err error) {
	m.State.Operation.InProgress = false
	m.State.Operation.Success = success
	m.State.Operation.Error = err
}

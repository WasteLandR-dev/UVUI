package app

import (
	"uvui/internal/services"
	"uvui/internal/types"
	"uvui/internal/ui/panels"
)

// Model represents the application state and dependencies
type Model struct {
	State           *types.AppState
	UVInstaller     services.UVInstallerInterface
	PythonManager   services.PythonManagerInterface
	ProjectManager  services.ProjectManagerInterface
	CommandExecutor services.CommandExecutorInterface
}

// NewModel creates a new application model
func NewModel() *Model {
	executor := services.NewCommandExecutor()

	state := &types.AppState{
		ActivePanel: types.StatusPanel,
		Panels: []types.Panel{
			types.StatusPanel,
			types.PythonPanel,
			types.ProjectPanel,
			types.EnvironmentPanel,
		},
		PythonVersions: types.PythonVersions{
			Available: []types.PythonVersion{},
			Installed: []types.PythonVersion{},
			Selected:  0,
			Loading:   false,
		},
		ProjectState: types.ProjectState{
			Status:         nil,
			Dependencies:   []types.ProjectDependency{},
			DependencyTree: nil,
			Selected:       0,
			Loading:        false,
			ShowTree:       false,
		},
		Messages: []string{},
	}

	return &Model{
		State:           state,
		UVInstaller:     services.NewUVInstaller(executor),
		PythonManager:   services.NewPythonManager(executor),
		ProjectManager:  services.NewProjectManager(executor),
		CommandExecutor: executor,
	}
}

// AddMessage adds a message to the message list
func (m *Model) AddMessage(msg string) {
	m.State.Messages = append(m.State.Messages, msg)
	// Keep only last 10 messages
	if len(m.State.Messages) > 10 {
		m.State.Messages = m.State.Messages[1:]
	}
}

// GetSelectedPythonVersion returns the currently selected Python version
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

// ValidateAndFixSelection ensures the selection index is within bounds
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

// GetMergedPythonVersions returns the merged list of available and installed versions
func (m *Model) GetMergedPythonVersions() []types.PythonVersion {
	// Import the panels package to use MergePythonVersions
	// This is a bit of a circular dependency, but necessary for consistency
	return panels.MergePythonVersions(m.State.PythonVersions.Available, m.State.PythonVersions.Installed)
}

// GetMergedPythonVersionsCount returns the count of merged Python versions
func (m *Model) GetMergedPythonVersionsCount() int {
	return len(m.GetMergedPythonVersions())
}

// UpdatePythonVersions updates the Python versions in the state
func (m *Model) UpdatePythonVersions(available, installed []types.PythonVersion) {
	m.State.PythonVersions.Available = available
	m.State.PythonVersions.Installed = installed
	m.State.PythonVersions.Loading = false

	// Validate and fix selection after updating versions
	m.ValidateAndFixSelection()
}

// UpdateProjectStatus updates the project status in the state
func (m *Model) UpdateProjectStatus(status *types.ProjectStatus) {
	m.State.ProjectState.Status = status
	m.State.ProjectState.Loading = false
}

// UpdateProjectDependencies updates project dependencies and tree
func (m *Model) UpdateProjectDependencies(deps []types.ProjectDependency, tree *types.DependencyTree) {
	m.State.ProjectState.Dependencies = deps
	m.State.ProjectState.DependencyTree = tree
}

// ToggleTreeView toggles between dependency list and tree view
func (m *Model) ToggleTreeView() {
	m.State.ProjectState.ShowTree = !m.State.ProjectState.ShowTree
}

// SetProjectLoading sets the project loading state
func (m *Model) SetProjectLoading(loading bool) {
	m.State.ProjectState.Loading = loading
}

// SetOperation sets the current operation status
func (m *Model) SetOperation(operation, target string, inProgress bool) {
	m.State.Operation = types.OperationStatus{
		InProgress: inProgress,
		Operation:  operation,
		Target:     target,
	}
}

// CompleteOperation completes the current operation
func (m *Model) CompleteOperation(success bool, err error) {
	m.State.Operation.InProgress = false
	m.State.Operation.Success = success
	m.State.Operation.Error = err
}

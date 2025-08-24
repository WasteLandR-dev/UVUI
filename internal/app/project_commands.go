package app

import (
	"uvui/internal/services"
	"uvui/internal/types"
	"uvui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// LoadProjectStatus loads the current project status
func LoadProjectStatus(projectManager services.ProjectManagerInterface) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		status, err := projectManager.GetProjectStatus()
		return ui.ProjectStatusLoadedMsg{
			Status: status,
			Error:  err,
		}
	})
}

// LoadProjectDependencies loads project dependencies and dependency tree
func LoadProjectDependencies(projectManager services.ProjectManagerInterface) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		deps, err1 := projectManager.GetProjectDependencies()
		tree, err2 := projectManager.GetDependencyTree()

		var err error
		if err1 != nil {
			err = err1
		} else if err2 != nil {
			err = err2
		}

		return ui.ProjectDependenciesLoadedMsg{
			Dependencies: deps,
			Tree:         tree,
			Error:        err,
		}
	})
}

// InitProject initializes a new project
func InitProject(projectManager services.ProjectManagerInterface, name string, options types.InitOptions) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := projectManager.InitProject(name, options)
		return ui.ProjectOperationMsg{
			Operation: "init",
			Success:   err == nil,
			Error:     err,
		}
	})
}

// SyncProject syncs project dependencies
func SyncProject(projectManager services.ProjectManagerInterface) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := projectManager.SyncProject()
		return ui.ProjectOperationMsg{
			Operation: "sync",
			Success:   err == nil,
			Error:     err,
		}
	})
}

// LockProject locks project dependencies
func LockProject(projectManager services.ProjectManagerInterface) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := projectManager.LockProject()
		return ui.ProjectOperationMsg{
			Operation: "lock",
			Success:   err == nil,
			Error:     err,
		}
	})
}

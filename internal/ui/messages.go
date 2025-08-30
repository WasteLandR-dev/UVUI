// Package ui provides the user interface for the application.
package ui

import "uvui/internal/types"

// UVInstalledMsg represents a message when UV installation is complete
type UVInstalledMsg struct {
	Success bool
	Version string
	Error   error
}

// PythonVersionsLoadedMsg represents loaded Python versions
type PythonVersionsLoadedMsg struct {
	Available []types.PythonVersion
	Installed []types.PythonVersion
	Error     error
}

// PythonOperationMsg represents a Python operation result
type PythonOperationMsg struct {
	Operation string
	Target    string
	Success   bool
	Error     error
}

// RefreshRequestMsg represents a request to refresh data
type RefreshRequestMsg struct {
	Panel types.Panel
}

// ErrorMsg represents an error message
type ErrorMsg struct {
	Error   error
	Context string
}

// StatusUpdateMsg represents a status update
type StatusUpdateMsg struct {
	Message string
	Type    string // "info", "warning", "error", "success"
}

// ProjectStatusLoadedMsg represents loaded project status
type ProjectStatusLoadedMsg struct {
	Status *types.ProjectStatus
	Error  error
}

// ProjectDependenciesLoadedMsg represents loaded project dependencies
type ProjectDependenciesLoadedMsg struct {
	Dependencies []types.ProjectDependency
	Tree         *types.DependencyTree
	Error        error
}

// ProjectOperationMsg represents a project operation result
type ProjectOperationMsg struct {
	Operation string // "init", "sync", "lock", "tree"
	Success   bool
	Error     error
}

// ProjectInitRequestMsg represents a project initialization request
type ProjectInitRequestMsg struct {
	Name    string
	Options types.InitOptions
}

// ToggleTreeViewMsg represents a request to toggle tree view
type ToggleTreeViewMsg struct{}

// RefreshProjectMsg represents a request to refresh project data
type RefreshProjectMsg struct{}

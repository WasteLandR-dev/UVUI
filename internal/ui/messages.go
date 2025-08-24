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

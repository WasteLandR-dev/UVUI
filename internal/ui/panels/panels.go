package panels

import "uvui/internal/types"

// AppState represents the overall application state.
type AppState struct {
	ActivePanel types.Panel
	Panels      []types.Panel
	Width       int
	Height      int
	types.UVStatus
	PythonVersions PythonVersions
	Installing     bool
	Messages       []string
	Operation      types.OperationStatus
	ProjectState   ProjectState
}

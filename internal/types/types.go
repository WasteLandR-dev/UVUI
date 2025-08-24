package types

// Panel represents a UI panel
type Panel int

const (
	StatusPanel Panel = iota
	PythonPanel
	ProjectPanel
	EnvironmentPanel
)

// UVStatus represents the status of UV installation
type UVStatus struct {
	Installed bool
	Version   string
	Path      string
}

// PythonVersion represents a Python version
type PythonVersion struct {
	Version   string
	Installed bool
	Current   bool
	Path      string
}

// PythonVersions represents the state of Python versions
type PythonVersions struct {
	Available []PythonVersion
	Installed []PythonVersion
	Selected  int
	Loading   bool
}

// OperationStatus represents the status of an ongoing operation
type OperationStatus struct {
	InProgress bool
	Operation  string
	Target     string
	Success    bool
	Error      error
}

// AppState represents the overall application state
type AppState struct {
	ActivePanel    Panel
	Panels         []Panel
	Width          int
	Height         int
	UVStatus       UVStatus
	PythonVersions PythonVersions
	Installing     bool
	Messages       []string
	Operation      OperationStatus
}

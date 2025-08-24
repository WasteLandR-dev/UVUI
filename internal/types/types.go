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
	ProjectState   ProjectState
}

// ProjectStatus represents the current project status
type ProjectStatus struct {
	IsProject     bool
	Name          string
	Path          string
	ConfigFile    string
	HasLockFile   bool
	LockFile      string
	PythonVersion string
	HasVirtualEnv bool
	VenvPath      string
}

// InitOptions represents options for project initialization
type InitOptions struct {
	App           bool
	Lib           bool
	PythonVersion string
}

// ProjectDependency represents a project dependency
type ProjectDependency struct {
	Name    string
	Version string
	Type    string // "main", "dev", etc.
}

// TreeNode represents a node in the dependency tree
type TreeNode struct {
	Name     string
	Version  string
	Level    int
	Children []TreeNode
}

// DependencyTree represents the project's dependency tree
type DependencyTree struct {
	Dependencies []TreeNode
}

// ProjectState represents the project panel state
type ProjectState struct {
	Status         *ProjectStatus
	Dependencies   []ProjectDependency
	DependencyTree *DependencyTree
	Selected       int
	Loading        bool
	ShowTree       bool
}

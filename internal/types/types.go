// Package types provides shared data types for the application.
package types

// Panel represents a UI panel.
type Panel int

const (
	// StatusPanel is the status panel.
	StatusPanel Panel = iota
	// PythonPanel is the python panel.
	PythonPanel
	// ProjectPanel is the project panel.
	ProjectPanel
	// EnvironmentPanel is the environment panel.
	EnvironmentPanel
)

// PythonVersion represents a Python version.
type PythonVersion struct {
	Version   string
	Installed bool
	Current   bool
	Path      string
}

// ProjectStatus represents the current project status.
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

// InitOptions represents options for project initialization.
type InitOptions struct {
	App           bool
	Lib           bool
	PythonVersion string
}

// ProjectDependency represents a project dependency.
type ProjectDependency struct {
	Name    string
	Version string
	Type    string // "main", "dev", etc.
}

// TreeNode represents a node in the dependency tree.
type TreeNode struct {
	Name     string
	Version  string
	Level    int
	Children []TreeNode
}

// DependencyTree represents the project's dependency tree.
type DependencyTree struct {
	Dependencies []TreeNode
}

// UVStatus represents the status of UV installation.
type UVStatus struct {
	Installed bool
	Version   string
	Path      string
}

// OperationStatus represents the status of an ongoing operation.
type OperationStatus struct {
	InProgress bool
	Operation  string
	Target     string
	Success    bool
	Error      error
}

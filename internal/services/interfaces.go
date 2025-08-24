package services

import "uvui/internal/types"

// CommandExecutorInterface defines the contract for command execution
type CommandExecutorInterface interface {
	Execute(command string, args ...string) ([]byte, error)
	IsUVAvailable() bool
}

// PythonManagerInterface defines the contract for Python version management
type PythonManagerInterface interface {
	ListAvailable() ([]types.PythonVersion, error)
	ListInstalled() ([]types.PythonVersion, error)
	Install(version string) error
	Uninstall(version string) error
	Pin(version string) error
	Find(version string) (*types.PythonVersion, error)
}

// UVInstallerInterface defines the contract for UV installation
type UVInstallerInterface interface {
	IsInstalled() (bool, string, error)
	Install() error
	GetInstallCommand() (string, error)
}

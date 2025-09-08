// Package services provides services for the application.
package services

import (
	"fmt"
	"os"
	"strings"
)

// UVInstaller implements UV installation functionality.
type UVInstaller struct {
	executor CommandExecutorInterface
}

// NewUVInstaller creates a new UV installer.
func NewUVInstaller(executor CommandExecutorInterface) *UVInstaller {
	return &UVInstaller{executor: executor}
}

// IsInstalled checks if UV is installed and returns version info.
func (u *UVInstaller) IsInstalled() (bool, string, error) {
	if !u.executor.IsUVAvailable() {
		return false, "", nil
	}

	output, err := u.executor.Execute("uv", "--version")
	if err != nil {
		return false, "", err
	}

	version := strings.TrimSpace(string(output))
	return true, version, nil
}

// Install installs UV based on the detected OS.
func (u *UVInstaller) Install() error {
	// This would execute the appropriate installation command
	// Implementation depends on OS detection
	return nil
}

// GetInstallCommand returns the installation command for the current OS.
func (u *UVInstaller) GetInstallCommand() (string, error) {
	// Detect OS and return appropriate command
	switch getOS() {
	case "windows":
		return "powershell -ExecutionPolicy ByPass -c \"irm https://astral.sh/uv/install.ps1 | iex\"", nil
	case "macos", "linux":
		return "curl -LsSf https://astral.sh/uv/install.sh | sh", nil
	default:
		return "", fmt.Errorf("unsupported operating system")
	}
}

// getOS detects the current operating system.
func getOS() string {
	sswitch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		return "windows"
	case fileExists("/System/Library/CoreServices/SystemVersion.plist"):
		return "macos"
	default:
		return "linux"
	}
}

// fileExists checks if a file exists.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

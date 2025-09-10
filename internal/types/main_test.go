package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanelConstants(t *testing.T) {
	// Test that panel constants are sequential and start from 0
	assert.Equal(t, 0, int(StatusPanel))
	assert.Equal(t, 1, int(PythonPanel))
	assert.Equal(t, 2, int(ProjectPanel))
	assert.Equal(t, 3, int(EnvironmentPanel))
}

func TestPythonVersion(t *testing.T) {
	version := PythonVersion{
		Version:   "3.12.0",
		Installed: true,
		Current:   false,
		Path:      "/usr/local/bin/python3.12",
	}

	assert.Equal(t, "3.12.0", version.Version)
	assert.True(t, version.Installed)
	assert.False(t, version.Current)
	assert.Equal(t, "/usr/local/bin/python3.12", version.Path)
}

func TestOperationStatus(t *testing.T) {
	status := OperationStatus{
		InProgress: true,
		Operation:  "install",
		Target:     "3.12.0",
		Success:    false,
		Error:      nil,
	}

	assert.True(t, status.InProgress)
	assert.Equal(t, "install", status.Operation)
	assert.Equal(t, "3.12.0", status.Target)
	assert.False(t, status.Success)
	assert.Nil(t, status.Error)
}

func TestUVStatus(t *testing.T) {
	status := UVStatus{
		Installed: true,
		Version:   "1.0.0",
		Path:      "/usr/local/bin/uv",
	}

	assert.True(t, status.Installed)
	assert.Equal(t, "1.0.0", status.Version)
	assert.Equal(t, "/usr/local/bin/uv", status.Path)
}

func TestPanelConversion(t *testing.T) {
	// Test converting panels to integers
	assert.Equal(t, 0, int(StatusPanel))
	assert.Equal(t, 1, int(PythonPanel))
	assert.Equal(t, 2, int(ProjectPanel))
	assert.Equal(t, 3, int(EnvironmentPanel))

	// Test converting integers back to panels
	assert.Equal(t, StatusPanel, Panel(0))
	assert.Equal(t, PythonPanel, Panel(1))
	assert.Equal(t, ProjectPanel, Panel(2))
	assert.Equal(t, EnvironmentPanel, Panel(3))
}

func TestPythonVersionComparison(t *testing.T) {
	version1 := PythonVersion{Version: "3.12.0", Installed: true}
	version2 := PythonVersion{Version: "3.11.0", Installed: false}
	version3 := PythonVersion{Version: "3.12.0", Installed: false}

	// Test that versions with same version string are equal
	assert.Equal(t, version1.Version, version3.Version)

	// Test that different versions are not equal
	assert.NotEqual(t, version1.Version, version2.Version)
}

func TestOperationStatusEmpty(t *testing.T) {
	var status OperationStatus

	assert.False(t, status.InProgress)
	assert.Equal(t, "", status.Operation)
	assert.Equal(t, "", status.Target)
	assert.False(t, status.Success)
	assert.Nil(t, status.Error)
}

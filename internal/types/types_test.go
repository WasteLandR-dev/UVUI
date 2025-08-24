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

func TestPythonVersions(t *testing.T) {
	versions := PythonVersions{
		Available: []PythonVersion{
			{Version: "3.12.0", Installed: false},
			{Version: "3.11.0", Installed: false},
		},
		Installed: []PythonVersion{
			{Version: "3.11.0", Installed: true, Current: true},
		},
		Selected: 0,
		Loading:  false,
	}

	assert.Len(t, versions.Available, 2)
	assert.Len(t, versions.Installed, 1)
	assert.Equal(t, 0, versions.Selected)
	assert.False(t, versions.Loading)
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

func TestAppState(t *testing.T) {
	state := AppState{
		ActivePanel: StatusPanel,
		Panels: []Panel{
			StatusPanel,
			PythonPanel,
			ProjectPanel,
			EnvironmentPanel,
		},
		Width:      100,
		Height:     50,
		Installing: false,
		Messages:   []string{},
	}

	assert.Equal(t, StatusPanel, state.ActivePanel)
	assert.Len(t, state.Panels, 4)
	assert.Equal(t, 100, state.Width)
	assert.Equal(t, 50, state.Height)
	assert.False(t, state.Installing)
	assert.Len(t, state.Messages, 0)
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

func TestEmptyAppState(t *testing.T) {
	// Test creating an empty app state
	var state AppState

	// Test default values
	assert.Equal(t, Panel(0), state.ActivePanel)
	assert.Nil(t, state.Panels)
	assert.Equal(t, 0, state.Width)
	assert.Equal(t, 0, state.Height)
	assert.False(t, state.Installing)
	assert.Nil(t, state.Messages)
}

func TestPythonVersionsEmpty(t *testing.T) {
	versions := PythonVersions{}

	assert.Nil(t, versions.Available)
	assert.Nil(t, versions.Installed)
	assert.Equal(t, 0, versions.Selected)
	assert.False(t, versions.Loading)
}

func TestOperationStatusEmpty(t *testing.T) {
	var status OperationStatus

	assert.False(t, status.InProgress)
	assert.Equal(t, "", status.Operation)
	assert.Equal(t, "", status.Target)
	assert.False(t, status.Success)
	assert.Nil(t, status.Error)
}

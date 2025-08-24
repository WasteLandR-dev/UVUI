package panels

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uvui/internal/types"
)

func TestRenderPythonPanel_UVNotInstalled(t *testing.T) {
	state := &types.AppState{
		UVStatus: types.UVStatus{Installed: false},
	}

	content := RenderPythonPanel(state)

	assert.Contains(t, content, "Python Version Management")
	assert.Contains(t, content, "UV must be installed first")
}

func TestRenderPythonPanel_Loading(t *testing.T) {
	state := &types.AppState{
		UVStatus: types.UVStatus{Installed: true},
		PythonVersions: types.PythonVersions{
			Loading: true,
		},
	}

	content := RenderPythonPanel(state)

	assert.Contains(t, content, "Python Version Management")
	assert.Contains(t, content, "Loading Python versions...")
}

func TestRenderPythonPanel_NoVersions(t *testing.T) {
	state := &types.AppState{
		UVStatus: types.UVStatus{Installed: true},
		PythonVersions: types.PythonVersions{
			Available: []types.PythonVersion{},
			Installed: []types.PythonVersion{},
		},
	}

	content := RenderPythonPanel(state)

	assert.Contains(t, content, "Python Version Management")
	assert.Contains(t, content, "No Python versions found")
	assert.Contains(t, content, "Press 'i' to refresh")
}

func TestRenderPythonPanel_WithVersions(t *testing.T) {
	state := &types.AppState{
		UVStatus: types.UVStatus{Installed: true},
		PythonVersions: types.PythonVersions{
			Available: []types.PythonVersion{
				{Version: "3.12.0", Installed: false},
				{Version: "3.11.0", Installed: false},
			},
			Installed: []types.PythonVersion{
				{Version: "3.11.0", Installed: true, Current: true, Path: "/usr/local/bin/python3.11"},
			},
			Selected: 0,
		},
	}

	content := RenderPythonPanel(state)

	assert.Contains(t, content, "Python Version Management")
	assert.Contains(t, content, "Python Versions (2 available)")
	assert.Contains(t, content, "3.12.0")
	assert.Contains(t, content, "3.11.0")
	assert.Contains(t, content, "(current)")
	assert.Contains(t, content, "(current)")
}

func TestRenderPythonPanel_OperationInProgress(t *testing.T) {
	state := &types.AppState{
		UVStatus: types.UVStatus{Installed: true},
		PythonVersions: types.PythonVersions{
			Available: []types.PythonVersion{
				{Version: "3.12.0", Installed: false},
			},
			Installed: []types.PythonVersion{},
		},
		Operation: types.OperationStatus{
			InProgress: true,
			Operation:  "install",
			Target:     "3.12.0",
		},
	}

	content := RenderPythonPanel(state)

	assert.Contains(t, content, "Install Python 3.12.0...")
}

func TestRenderVersionLine_Selected(t *testing.T) {
	version := types.PythonVersion{
		Version:   "3.12.0",
		Installed: false,
		Current:   false,
	}

	line := renderVersionLine(version, true)

	assert.Contains(t, line, "> ")
	assert.Contains(t, line, "3.12.0")
}

func TestRenderVersionLine_NotSelected(t *testing.T) {
	version := types.PythonVersion{
		Version:   "3.12.0",
		Installed: false,
		Current:   false,
	}

	line := renderVersionLine(version, false)

	assert.Contains(t, line, "  ")
	assert.Contains(t, line, "3.12.0")
}

func TestRenderVersionLine_Current(t *testing.T) {
	version := types.PythonVersion{
		Version:   "3.11.0",
		Installed: true,
		Current:   true,
		Path:      "/usr/local/bin/python3.11",
	}

	line := renderVersionLine(version, false)

	assert.Contains(t, line, "3.11.0")
	assert.Contains(t, line, "(current)")
	assert.Contains(t, line, "/usr/local/bin/python3.11")
}

func TestRenderVersionLine_Installed(t *testing.T) {
	version := types.PythonVersion{
		Version:   "3.10.0",
		Installed: true,
		Current:   false,
		Path:      "/usr/local/bin/python3.10",
	}

	line := renderVersionLine(version, false)

	assert.Contains(t, line, "3.10.0")
	assert.Contains(t, line, "✓")
	assert.Contains(t, line, "/usr/local/bin/python3.10")
}

func TestMergePythonVersions_Empty(t *testing.T) {
	available := []types.PythonVersion{}
	installed := []types.PythonVersion{}

	merged := MergePythonVersions(available, installed)

	assert.Len(t, merged, 0)
}

func TestMergePythonVersions_OnlyAvailable(t *testing.T) {
	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	installed := []types.PythonVersion{}

	merged := MergePythonVersions(available, installed)

	assert.Len(t, merged, 2)
	assert.Equal(t, "3.12.0", merged[0].Version)
	assert.Equal(t, "3.11.0", merged[1].Version)
}

func TestMergePythonVersions_OnlyInstalled(t *testing.T) {
	available := []types.PythonVersion{}
	installed := []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true},
		{Version: "3.10.0", Installed: true, Current: false},
	}

	merged := MergePythonVersions(available, installed)

	assert.Len(t, merged, 2)
	assert.Equal(t, "3.11.0", merged[0].Version)
	assert.Equal(t, "3.10.0", merged[1].Version)
}

func TestMergePythonVersions_Overlap(t *testing.T) {
	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	installed := []types.PythonVersion{
		{Version: "3.11.0", Installed: true, Current: true, Path: "/usr/local/bin/python3.11"},
	}

	merged := MergePythonVersions(available, installed)

	assert.Len(t, merged, 2)

	// Find the installed version
	var installedVersion *types.PythonVersion
	for _, v := range merged {
		if v.Version == "3.11.0" {
			installedVersion = &v
			break
		}
	}

	assert.NotNil(t, installedVersion)
	assert.True(t, installedVersion.Installed)
	assert.True(t, installedVersion.Current)
	assert.Equal(t, "/usr/local/bin/python3.11", installedVersion.Path)
}

func TestMergePythonVersions_Sorting(t *testing.T) {
	available := []types.PythonVersion{
		{Version: "3.10.0", Installed: false},
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	installed := []types.PythonVersion{}

	merged := MergePythonVersions(available, installed)

	assert.Len(t, merged, 3)
	// Should be sorted in descending order (highest first)
	assert.Equal(t, "3.12.0", merged[0].Version)
	assert.Equal(t, "3.11.0", merged[1].Version)
	assert.Equal(t, "3.10.0", merged[2].Version)
}

func TestGetPythonPanelHelp(t *testing.T) {
	help := GetPythonPanelHelp()

	assert.Contains(t, help, "↑↓: Navigate")
	assert.Contains(t, help, "Enter: Install")
	assert.Contains(t, help, "d: Delete")
	assert.Contains(t, help, "p: Pin")
	assert.Contains(t, help, "i: Refresh")
}

func TestCompareVersions(t *testing.T) {
	// Test basic version comparison using the exported function
	// Note: This tests the actual implementation in the panels package
	available := []types.PythonVersion{
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.0", Installed: false},
	}
	installed := []types.PythonVersion{}

	merged := MergePythonVersions(available, installed)

	// Should be sorted in descending order (highest first)
	assert.Equal(t, "3.12.0", merged[0].Version)
	assert.Equal(t, "3.11.0", merged[1].Version)
}

func TestParseVersionPart(t *testing.T) {
	// Test numeric parts
	assert.Equal(t, 12, parseVersionPart("12"))
	assert.Equal(t, 1, parseVersionPart("1"))
	assert.Equal(t, 0, parseVersionPart("0"))

	// Test with non-numeric suffix
	assert.Equal(t, 12, parseVersionPart("12alpha"))
	assert.Equal(t, 1, parseVersionPart("1beta"))

	// Test non-numeric
	assert.Equal(t, 0, parseVersionPart("alpha"))
	assert.Equal(t, 0, parseVersionPart(""))
}

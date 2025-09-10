package services

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"uvui/internal/types"
)

func TestNewProjectManager(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewProjectManager(executor)
	if pm == nil {
		t.Error("NewProjectManager() should not return nil")
	}
}

func TestGetProjectStatus_NoProject(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewProjectManager(executor)

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "uvui-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Change the current working directory to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	status, err := pm.GetProjectStatus()
	if err != nil {
		t.Errorf("GetProjectStatus() error = %v, wantErr %v", err, false)
	}

	if status.IsProject {
		t.Error("GetProjectStatus() IsProject = true, want false")
	}
}

func TestGetProjectStatus_Project(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewProjectManager(executor)

	// Create a temporary directory and a pyproject.toml file
	tmpDir, err := os.MkdirTemp("", "uvui-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	if _, err := os.Create(pyprojectPath); err != nil {
		t.Fatal(err)
	}

	// Change the current working directory to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	status, err := pm.GetProjectStatus()
	if err != nil {
		t.Errorf("GetProjectStatus() error = %v, wantErr %v", err, false)
	}

	if !status.IsProject {
		t.Error("GetProjectStatus() IsProject = false, want true")
	}
}

func TestInitProject(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "init" {
				return nil, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewProjectManager(executor)

	projectName, err := pm.InitProject("test-project", types.InitOptions{})
	if err != nil {
		t.Errorf("InitProject() error = %v, wantErr %v", err, false)
	}

	if projectName != "test-project" {
		t.Errorf("InitProject() projectName = %q, want %q", projectName, "test-project")
	}
}

func TestSyncProject(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "sync" {
				return nil, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewProjectManager(executor)

	err := pm.SyncProject()
	if err != nil {
		t.Errorf("SyncProject() error = %v, wantErr %v", err, false)
	}
}

func TestLockProject(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "lock" {
				return nil, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewProjectManager(executor)

	err := pm.LockProject()
	if err != nil {
		t.Errorf("LockProject() error = %v, wantErr %v", err, false)
	}
}

func TestGetDependencyTree(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "tree" {
				return []byte("requests==2.31.0\n"), nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewProjectManager(executor)

	tree, err := pm.GetDependencyTree()
	if err != nil {
		t.Errorf("GetDependencyTree() error = %v, wantErr %v", err, false)
	}

	expectedTree := &types.DependencyTree{
		Dependencies: []types.TreeNode{
			{Name: "requests", Version: "2.31.0", Level: 0},
		},
	}

	if !reflect.DeepEqual(tree, expectedTree) {
		t.Errorf("GetDependencyTree() tree = %v, want %v", tree, expectedTree)
	}
}

func TestGetProjectDependencies(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewProjectManager(executor)

	// Create a temporary directory and a pyproject.toml file
	tmpDir, err := os.MkdirTemp("", "uvui-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	if _, err := os.Create(pyprojectPath); err != nil {
		t.Fatal(err)
	}

	// Change the current working directory to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	deps, err := pm.GetProjectDependencies()
	if err != nil {
		t.Errorf("GetProjectDependencies() error = %v, wantErr %v", err, false)
	}

	if len(deps) == 0 {
		t.Error("GetProjectDependencies() deps is empty, want not empty")
	}
}

func TestInitProject_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewProjectManager(executor)

	_, err := pm.InitProject("test-project", types.InitOptions{})
	if err == nil {
		t.Error("InitProject() error = nil, wantErr true")
	}
}

func TestSyncProject_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewProjectManager(executor)

	err := pm.SyncProject()
	if err == nil {
		t.Error("SyncProject() error = nil, wantErr true")
	}
}

func TestLockProject_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewProjectManager(executor)

	err := pm.LockProject()
	if err == nil {
		t.Error("LockProject() error = nil, wantErr true")
	}
}

func TestGetDependencyTree_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewProjectManager(executor)

	_, err := pm.GetDependencyTree()
	if err == nil {
		t.Error("GetDependencyTree() error = nil, wantErr true")
	}
}

func TestGetProjectDependencies_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		IsUVAvailableFunc: func() bool {
			return false
		},
	}
	pm := NewProjectManager(executor)

	_, err := pm.GetProjectDependencies()
	if err == nil {
		t.Error("GetProjectDependencies() error = nil, wantErr true")
	}
}

func TestGetProjectStatus_Error(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewProjectManager(executor)

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "uvui-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Change the current working directory to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	// Create a pyproject.toml file
	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	if _, err := os.Create(pyprojectPath); err != nil {
		t.Fatal(err)
	}

	// Make the directory unreadable
	if err := os.Chmod(tmpDir, 0000); err != nil {
		t.Fatal(err)
	}

	_, err = pm.GetProjectStatus()
	if err == nil {
		t.Error("GetProjectStatus() error = nil, wantErr true")
	}

	// Restore permissions so the directory can be cleaned up
	_ = os.Chmod(tmpDir, 0755)
}

func TestGetProjectStatus_Project_NoLockOrVersion(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewProjectManager(executor)

	// Create a temporary directory and a pyproject.toml file
	tmpDir, err := os.MkdirTemp("", "uvui-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	if _, err := os.Create(pyprojectPath); err != nil {
		t.Fatal(err)
	}

	// Change the current working directory to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	status, err := pm.GetProjectStatus()
	if err != nil {
		t.Errorf("GetProjectStatus() error = %v, wantErr %v", err, false)
	}

	if !status.IsProject {
		t.Error("GetProjectStatus() IsProject = false, want true")
	}

	if status.HasLockFile {
		t.Error("GetProjectStatus() HasLockFile = true, want false")
	}

	if status.PythonVersion != "" {
		t.Errorf("GetProjectStatus() PythonVersion = %q, want %q", status.PythonVersion, "")
	}
}

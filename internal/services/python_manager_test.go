package services

import (
	"fmt"
	"reflect"
	"testing"

	"uvui/internal/types"
)

func TestNewPythonManager(t *testing.T) {
	executor := &mockCommandExecutor{}
	pm := NewPythonManager(executor)
	if pm == nil {
		t.Error("NewPythonManager() should not return nil")
	}
}

func TestPythonManager_ListAvailable(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "python" && args[1] == "list" {
				return []byte("3.12.1\n3.12.0\n"), nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewPythonManager(executor)

	versions, err := pm.ListAvailable()
	if err != nil {
		t.Errorf("ListAvailable() error = %v, wantErr %v", err, false)
	}

	expectedVersions := []types.PythonVersion{
		{Version: "3.12.1", Installed: false, Current: false},
		{Version: "3.12.0", Installed: false, Current: false},
	}

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("ListAvailable() versions = %v, want %v", versions, expectedVersions)
	}
}

func TestPythonManager_ListInstalled(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "python" && args[1] == "list" {
				return []byte("3.11.6 * /path/to/python3.11\n3.10.12 /path/to/python3.10"), nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewPythonManager(executor)

	versions, err := pm.ListInstalled()
	if err != nil {
		t.Errorf("ListInstalled() error = %v, wantErr %v", err, false)
	}

	expectedVersions := []types.PythonVersion{
		{Version: "3.11.6", Installed: true, Current: true, Path: "/path/to/python3.11"},
		{Version: "3.10.12", Installed: true, Current: false, Path: "/path/to/python3.10"},
	}

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("ListInstalled() versions = %v, want %v", versions, expectedVersions)
	}
}

func TestPythonManager_Install(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "python" && args[1] == "install" && args[2] == "3.12.1" {
				return nil, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewPythonManager(executor)

	err := pm.Install("3.12.1")
	if err != nil {
		t.Errorf("Install() error = %v, wantErr %v", err, false)
	}
}

func TestPythonManager_Uninstall(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "python" && args[1] == "uninstall" && args[2] == "3.12.1" {
				return nil, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewPythonManager(executor)

	err := pm.Uninstall("3.12.1")
	if err != nil {
		t.Errorf("Uninstall() error = %v, wantErr %v", err, false)
	}
}

func TestPythonManager_Pin(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "python" && args[1] == "pin" && args[2] == "3.12.1" {
				return nil, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewPythonManager(executor)

	err := pm.Pin("3.12.1")
	if err != nil {
		t.Errorf("Pin() error = %v, wantErr %v", err, false)
	}
}

func TestPythonManager_Find(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "python" && args[1] == "find" && args[2] == "3.12.1" {
				return []byte("3.12.1 /path/to/python3.12"), nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	pm := NewPythonManager(executor)

	version, err := pm.Find("3.12.1")
	if err != nil {
		t.Errorf("Find() error = %v, wantErr %v", err, false)
	}

	expectedVersion := &types.PythonVersion{
		Version:   "3.12.1",
		Installed: true,
		Path:      "/path/to/python3.12",
	}

	if !reflect.DeepEqual(version, expectedVersion) {
		t.Errorf("Find() version = %v, want %v", version, expectedVersion)
	}
}

func TestPythonManager_ListAvailable_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewPythonManager(executor)

	_, err := pm.ListAvailable()
	if err == nil {
		t.Error("ListAvailable() error = nil, wantErr true")
	}
}

func TestPythonManager_ListInstalled_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewPythonManager(executor)

	_, err := pm.ListInstalled()
	if err == nil {
		t.Error("ListInstalled() error = nil, wantErr true")
	}
}

func TestPythonManager_Install_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewPythonManager(executor)

	err := pm.Install("3.12.1")
	if err == nil {
		t.Error("Install() error = nil, wantErr true")
	}
}

func TestPythonManager_Uninstall_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewPythonManager(executor)

	err := pm.Uninstall("3.12.1")
	if err == nil {
		t.Error("Uninstall() error = nil, wantErr true")
	}
}

func TestPythonManager_Pin_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewPythonManager(executor)

	err := pm.Pin("3.12.1")
	if err == nil {
		t.Error("Pin() error = nil, wantErr true")
	}
}

func TestPythonManager_Find_Error(t *testing.T) {
	executor := &mockCommandExecutor{
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	pm := NewPythonManager(executor)

	_, err := pm.Find("3.12.1")
	if err == nil {
		t.Error("Find() error = nil, wantErr true")
	}
}

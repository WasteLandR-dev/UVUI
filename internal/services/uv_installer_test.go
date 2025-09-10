package services

import (
	"fmt"
	"testing"
)

func TestNewUVInstaller(t *testing.T) {
	executor := &mockCommandExecutor{}
	installer := NewUVInstaller(executor)
	if installer == nil {
		t.Error("NewUVInstaller() should not return nil")
	}
}

func TestUVInstaller_IsInstalled_NotInstalled(t *testing.T) {
	executor := &mockCommandExecutor{
		IsUVAvailableFunc: func() bool {
			return false
		},
	}
	installer := NewUVInstaller(executor)

	installed, version, err := installer.IsInstalled()
	if err != nil {
		t.Errorf("IsInstalled() error = %v, wantErr %v", err, false)
	}

	if installed {
		t.Error("IsInstalled() installed = true, want false")
	}

	if version != "" {
		t.Errorf("IsInstalled() version = %q, want %q", version, "")
	}
}

func TestUVInstaller_IsInstalled_Installed(t *testing.T) {
	executor := &mockCommandExecutor{
		IsUVAvailableFunc: func() bool {
			return true
		},
		ExecuteFunc: func(command string, args ...string) ([]byte, error) {
			if command == "uv" && args[0] == "--version" {
				return []byte("uv 0.1.0"), nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", command, args)
		},
	}
	installer := NewUVInstaller(executor)

	installed, version, err := installer.IsInstalled()
	if err != nil {
		t.Errorf("IsInstalled() error = %v, wantErr %v", err, false)
	}

	if !installed {
		t.Error("IsInstalled() installed = false, want true")
	}

	if version != "uv 0.1.0" {
		t.Errorf("IsInstalled() version = %q, want %q", version, "uv 0.1.0")
	}
}

func TestUVInstaller_Install(t *testing.T) {
	executor := &mockCommandExecutor{}
	installer := NewUVInstaller(executor)

	err := installer.Install()
	if err != nil {
		t.Errorf("Install() error = %v, wantErr %v", err, false)
	}
}

func TestUVInstaller_GetInstallCommand(t *testing.T) {
	executor := &mockCommandExecutor{}
	installer := NewUVInstaller(executor)

	_, err := installer.GetInstallCommand()
	if err != nil {
		t.Errorf("GetInstallCommand() error = %v, wantErr %v", err, false)
	}
}

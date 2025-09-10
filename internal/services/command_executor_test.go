package services

import (
	"testing"
)

func TestNewCommandExecutor(t *testing.T) {
	executor := NewCommandExecutor()
	if executor == nil {
		t.Error("NewCommandExecutor() should not return nil")
	}
}

func TestCommandExecutor_Execute(t *testing.T) {
	var executor CommandExecutorInterface = NewCommandExecutor()
	output, err := executor.Execute("echo", "hello")
	if err != nil {
		t.Errorf("Execute() error = %v, wantErr %v", err, false)
	}
	if string(output) != "hello\n" {
		t.Errorf("Execute() output = %q, want %q", string(output), "hello\n")
	}
}

func TestCommandExecutor_Execute_Error(t *testing.T) {
	var executor CommandExecutorInterface = NewCommandExecutor()
	_, err := executor.Execute("non-existent-command")
	if err == nil {
		t.Error("Execute() error = nil, wantErr true")
	}
}

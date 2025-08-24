package services

import (
	"os/exec"
)

// CommandExecutor implements command execution functionality
type CommandExecutor struct{}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

// Execute runs a command and returns its output
func (c *CommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	return cmd.Output()
}

// IsUVAvailable checks if UV is available in PATH
func (c *CommandExecutor) IsUVAvailable() bool {
	_, err := exec.LookPath("uv")
	return err == nil
}

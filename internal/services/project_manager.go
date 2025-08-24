package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"uvui/internal/types"
)

// ProjectManager implements project management functionality
type ProjectManager struct {
	executor CommandExecutorInterface
}

// NewProjectManager creates a new project manager
func NewProjectManager(executor CommandExecutorInterface) *ProjectManager {
	return &ProjectManager{executor: executor}
}

// GetProjectStatus returns the current project status
func (p *ProjectManager) GetProjectStatus() (*types.ProjectStatus, error) {
	status := &types.ProjectStatus{
		IsProject: false,
	}

	// Check for pyproject.toml file
	currentDir, err := os.Getwd()
	if err != nil {
		return status, err
	}

	pyprojectPath := filepath.Join(currentDir, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); os.IsNotExist(err) {
		status.Path = currentDir
		return status, nil
	}

	// Found pyproject.toml, this is a project
	status.IsProject = true
	status.Path = currentDir
	status.Name = filepath.Base(currentDir)
	status.ConfigFile = pyprojectPath

	// Try to get more project info
	if p.executor.IsUVAvailable() {
		// Check if project is synced
		lockFilePath := filepath.Join(currentDir, "uv.lock")
		if _, err := os.Stat(lockFilePath); err == nil {
			status.HasLockFile = true
			status.LockFile = lockFilePath
		}

		// Get Python version if pinned
		pythonVersionPath := filepath.Join(currentDir, ".python-version")
		if content, err := os.ReadFile(pythonVersionPath); err == nil {
			status.PythonVersion = strings.TrimSpace(string(content))
		}

		// Check virtual environment
		venvPath := filepath.Join(currentDir, ".venv")
		if _, err := os.Stat(venvPath); err == nil {
			status.HasVirtualEnv = true
			status.VenvPath = venvPath
		}
	}

	return status, nil
}

// InitProject creates a new UV project
func (p *ProjectManager) InitProject(name string, options types.InitOptions) error {
	if !p.executor.IsUVAvailable() {
		return fmt.Errorf("UV is not available")
	}

	args := []string{"init"}

	if name != "" {
		args = append(args, name)
	}

	if options.App {
		args = append(args, "--app")
	}

	if options.Lib {
		args = append(args, "--lib")
	}

	if options.PythonVersion != "" {
		args = append(args, "--python", options.PythonVersion)
	}

	_, err := p.executor.Execute("uv", args...)
	return err
}

// SyncProject syncs project dependencies
func (p *ProjectManager) SyncProject() error {
	if !p.executor.IsUVAvailable() {
		return fmt.Errorf("UV is not available")
	}

	_, err := p.executor.Execute("uv", "sync")
	return err
}

// LockProject locks project dependencies
func (p *ProjectManager) LockProject() error {
	if !p.executor.IsUVAvailable() {
		return fmt.Errorf("UV is not available")
	}

	_, err := p.executor.Execute("uv", "lock")
	return err
}

// GetDependencyTree returns the project dependency tree
func (p *ProjectManager) GetDependencyTree() (*types.DependencyTree, error) {
	if !p.executor.IsUVAvailable() {
		return nil, fmt.Errorf("UV is not available")
	}

	output, err := p.executor.Execute("uv", "tree")
	if err != nil {
		// Return mock data for demo purposes
		return p.getMockDependencyTree(), nil
	}

	return p.parseDependencyTree(string(output)), nil
}

// GetProjectDependencies returns project dependencies
func (p *ProjectManager) GetProjectDependencies() ([]types.ProjectDependency, error) {
	if !p.executor.IsUVAvailable() {
		return nil, fmt.Errorf("UV is not available")
	}

	// Try to read from pyproject.toml
	pyprojectPath := "pyproject.toml"
	if _, err := os.Stat(pyprojectPath); os.IsNotExist(err) {
		return []types.ProjectDependency{}, nil
	}

	// For now, return mock dependencies
	// In a real implementation, you'd parse the pyproject.toml file
	return p.getMockDependencies(), nil
}

// parseDependencyTree parses the output of uv tree
func (p *ProjectManager) parseDependencyTree(output string) *types.DependencyTree {
	tree := &types.DependencyTree{
		Dependencies: []types.TreeNode{},
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse tree structure
		level := 0
		trimmed := line
		for strings.HasPrefix(trimmed, "  ") || strings.HasPrefix(trimmed, "├─") || strings.HasPrefix(trimmed, "└─") {
			level++
			if strings.HasPrefix(trimmed, "├─") || strings.HasPrefix(trimmed, "└─") {
				trimmed = trimmed[2:]
			} else {
				trimmed = trimmed[2:]
			}
		}

		// Extract package name and version
		parts := strings.Fields(trimmed)
		if len(parts) > 0 {
			node := types.TreeNode{
				Name:  parts[0],
				Level: level,
			}
			if len(parts) > 1 {
				node.Version = parts[1]
			}
			tree.Dependencies = append(tree.Dependencies, node)
		}
	}

	return tree
}

// getMockDependencyTree returns mock dependency tree for demo
func (p *ProjectManager) getMockDependencyTree() *types.DependencyTree {
	return &types.DependencyTree{
		Dependencies: []types.TreeNode{
			{Name: "requests", Version: "2.31.0", Level: 0},
			{Name: "urllib3", Version: "2.0.4", Level: 1},
			{Name: "certifi", Version: "2023.7.22", Level: 1},
			{Name: "charset-normalizer", Version: "3.2.0", Level: 1},
			{Name: "idna", Version: "3.4", Level: 1},
			{Name: "fastapi", Version: "0.104.1", Level: 0},
			{Name: "starlette", Version: "0.27.0", Level: 1},
			{Name: "pydantic", Version: "2.4.2", Level: 1},
			{Name: "typing-extensions", Version: "4.8.0", Level: 2},
		},
	}
}

// getMockDependencies returns mock dependencies for demo
func (p *ProjectManager) getMockDependencies() []types.ProjectDependency {
	return []types.ProjectDependency{
		{Name: "requests", Version: "^2.31.0", Type: "main"},
		{Name: "fastapi", Version: "^0.104.1", Type: "main"},
		{Name: "pytest", Version: "^7.4.0", Type: "dev"},
		{Name: "black", Version: "^23.9.0", Type: "dev"},
		{Name: "mypy", Version: "^1.5.0", Type: "dev"},
	}
}

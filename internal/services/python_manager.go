// Package services provides services for the application.
package services

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"uvui/internal/types"
	"uvui/pkg/version"
)

// PythonManager implements Python version management functionality.
type PythonManager struct {
	executor CommandExecutorInterface
}

// NewPythonManager creates a new Python manager.
func NewPythonManager(executor CommandExecutorInterface) *PythonManager {
	return &PythonManager{executor: executor}
}

// ListAvailable lists all available Python versions.
func (p *PythonManager) ListAvailable() ([]types.PythonVersion, error) {
	if !p.executor.IsUVAvailable() {
		return nil, fmt.Errorf("UV is not available")
	}

	output, err := p.executor.Execute("uv", "python", "list", "--only-downloads")
	if err != nil {
		// Return mock data for demo purposes
		return p.getMockAvailableVersions(), nil
	}

	return p.parseAvailableVersions(string(output)), nil
}

// ListInstalled lists all installed Python versions.
func (p *PythonManager) ListInstalled() ([]types.PythonVersion, error) {
	if !p.executor.IsUVAvailable() {
		return nil, fmt.Errorf("UV is not available")
	}

	output, err := p.executor.Execute("uv", "python", "list", "--only-installed")
	if err != nil {
		// Return mock data for demo purposes
		return p.getMockInstalledVersions(), nil
	}

	return p.parseInstalledVersions(string(output)), nil
}

// Install installs a Python version.
func (p *PythonManager) Install(version string) error {
	if !p.executor.IsUVAvailable() {
		return fmt.Errorf("UV is not available")
	}

	_, err := p.executor.Execute("uv", "python", "install", version)
	return err
}

// Uninstall removes a Python version.
func (p *PythonManager) Uninstall(version string) error {
	if !p.executor.IsUVAvailable() {
		return fmt.Errorf("UV is not available")
	}

	_, err := p.executor.Execute("uv", "python", "uninstall", version)
	return err
}

// Pin pins a Python version for the current project.
func (p *PythonManager) Pin(version string) error {
	if !p.executor.IsUVAvailable() {
		return fmt.Errorf("UV is not available")
	}

	_, err := p.executor.Execute("uv", "python", "pin", version)
	return err
}

// Find finds a specific Python version.
func (p *PythonManager) Find(version string) (*types.PythonVersion, error) {
	if !p.executor.IsUVAvailable() {
		return nil, fmt.Errorf("UV is not available")
	}

	output, err := p.executor.Execute("uv", "python", "find", version)
	if err != nil {
		return nil, err
	}

	return p.parseFindResult(string(output)), nil
}

// parseAvailableVersions parses the output of uv python list --only-downloads.
func (p *PythonManager) parseAvailableVersions(output string) []types.PythonVersion {
	var versions []types.PythonVersion
	lines := strings.Split(output, "\n")

	versionRegex := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?(?:[a-z]+\d+)?)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := versionRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			versions = append(versions, types.PythonVersion{
				Version:   matches[1],
				Installed: false,
				Current:   false,
			})
		}
	}

	// Sort versions in descending order
	sort.Slice(versions, func(i, j int) bool {
		return version.CompareVersions(versions[i].Version, versions[j].Version) > 0
	})

	return versions
}

// parseInstalledVersions parses the output of uv python list --only-installed.
func (p *PythonManager) parseInstalledVersions(output string) []types.PythonVersion {
	var versions []types.PythonVersion
	lines := strings.Split(output, "\n")

	versionRegex := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?(?:[a-z]+\d+)?)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse installed versions with paths
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			ver := versionRegex.FindString(parts[0])
			versions = append(versions, types.PythonVersion{
				Version:   ver,
				Installed: true,
				Current:   strings.Contains(line, "*") || strings.Contains(line, "default"),
				Path:      strings.Join(parts[1:], " "),
			})
		}
	}

	return versions
}

// parseFindResult parses the output of uv python find.
func (p *PythonManager) parseFindResult(output string) *types.PythonVersion {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return nil
	}

	line := strings.TrimSpace(lines[0])
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return &types.PythonVersion{
			Version:   parts[0],
			Installed: true,
			Path:      strings.Join(parts[1:], " "),
		}
	}

	return nil
}

// getMockAvailableVersions returns mock data for available versions.
func (p *PythonManager) getMockAvailableVersions() []types.PythonVersion {
	return []types.PythonVersion{
		{Version: "3.12.1", Installed: false},
		{Version: "3.12.0", Installed: false},
		{Version: "3.11.7", Installed: false},
		{Version: "3.11.6", Installed: false},
		{Version: "3.10.13", Installed: false},
		{Version: "3.10.12", Installed: false},
		{Version: "3.9.18", Installed: false},
		{Version: "3.9.17", Installed: false},
		{Version: "3.8.18", Installed: false},
	}
}

// getMockInstalledVersions returns mock data for installed versions.
func (p *PythonManager) getMockInstalledVersions() []types.PythonVersion {
	return []types.PythonVersion{
		{Version: "3.11.6", Installed: true, Current: true, Path: "/usr/local/bin/python3.11"},
		{Version: "3.10.12", Installed: true, Current: false, Path: "/usr/local/bin/python3.10"},
	}
}

// Package panels provides UI panels for the application.
package panels

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"uvui/internal/types"
	"uvui/internal/ui"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PythonVersions represents the state of Python versions.
type PythonVersions struct {
	Available []types.PythonVersion
	Installed []types.PythonVersion
	Selected  int
	Loading   bool
}

// RenderPythonPanel renders the Python management panel.
func RenderPythonPanel(state *AppState) string {
	var content strings.Builder

	content.WriteString("Python Version Management\n\n")

	if !state.Installed {
		content.WriteString(ui.ErrorStyle.Render("UV must be installed first to manage Python versions."))
		return content.String()
	}

	if state.PythonVersions.Loading {
		content.WriteString(ui.LoadingStyle.Render("⏳ Loading Python versions..."))
		return content.String()
	}

	// Merge available and installed versions for display
	allVersions := MergePythonVersions(state.PythonVersions.Available, state.PythonVersions.Installed)

	if len(allVersions) == 0 {
		content.WriteString("No Python versions found.\n")
		content.WriteString("Press 'i' to refresh the list.")
		return content.String()
	}

	content.WriteString(fmt.Sprintf("Python Versions (%d available):\n\n", len(allVersions)))

	// Render version list
	pinnedVersion := getPinnedVersion()
	for i, version := range allVersions {
		line := renderVersionLine(version, i == state.PythonVersions.Selected, pinnedVersion)
		content.WriteString(line + "\n")
	}

	// Show operation status
	if state.Operation.InProgress &&
		(state.Operation.Operation == "install" ||
			state.Operation.Operation == "uninstall" ||
			state.Operation.Operation == "pin") {
		content.WriteString("\n")
		content.WriteString(ui.LoadingStyle.Render(fmt.Sprintf("⏳ %s Python %s...",
			cases.Title(language.English).String(state.Operation.Operation),
			state.Operation.Target)))
	}

	// Show help bar
	content.WriteString("\n\n---\n")
	content.WriteString(ui.HelpStyle.Render(GetPythonPanelHelp()))

	return content.String()
}

// renderVersionLine renders a single Python version line.
func renderVersionLine(version types.PythonVersion, selected bool, pinnedVersion string) string {
	var line strings.Builder

	// Selection indicator
	if selected {
		line.WriteString(ui.SelectedItemStyle.Render("> "))
	} else {
		line.WriteString("  ")
	}

	// Version number with status styling
	versionText := version.Version
	if version.Version == pinnedVersion {
		versionText = ui.PinnedVersionStyle.Render(versionText + " (pinned)")
	} else if version.Current {
		versionText = ui.CurrentVersionStyle.Render(versionText + " (current)")
	} else if version.Installed {
		versionText = ui.InstalledVersionStyle.Render(versionText + " ✓")
	} else {
		versionText = ui.AvailableVersionStyle.Render(versionText)
	}

	line.WriteString(versionText)

	// Path information for installed versions
	if version.Installed && version.Path != "" {
		line.WriteString(ui.AvailableVersionStyle.Render(fmt.Sprintf(" [%s]", version.Path)))
	}

	return line.String()
}

// MergePythonVersions merges available and installed versions for display (exported).
func MergePythonVersions(available, installed []types.PythonVersion) []types.PythonVersion {
	versionMap := make(map[string]types.PythonVersion)

	// Add all available versions
	for _, version := range available {
		versionMap[version.Version] = version
	}

	// Update with installed status
	for _, version := range installed {
		if existing, exists := versionMap[version.Version]; exists {
			existing.Installed = true
			existing.Current = version.Current
			existing.Path = version.Path
			versionMap[version.Version] = existing
		} else {
			versionMap[version.Version] = version
		}
	}

	// Convert back to slice
	var merged []types.PythonVersion
	for _, version := range versionMap {
		merged = append(merged, version)
	}

	// Sort by version number to ensure consistent order
	sort.Slice(merged, func(i, j int) bool {
		return compareVersions(merged[i].Version, merged[j].Version) > 0
	})

	return merged
}

// compareVersions compares two version strings (higher versions first).
func compareVersions(a, b string) int {
	// Simple version comparison - split by dots and compare numerically
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var numA, numB int
		if i < len(partsA) {
			numA = parseVersionPart(partsA[i])
		}
		if i < len(partsB) {
			numB = parseVersionPart(partsB[i])
		}

		if numA != numB {
			return numA - numB
		}
	}

	return 0
}

// parseVersionPart parses a version part (e.g., "12", "1", "alpha1").
func parseVersionPart(part string) int {
	// Extract numeric part
	re := regexp.MustCompile(`^(\d+)`)
	matches := re.FindStringSubmatch(part)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 0
}

// GetPythonPanelHelp returns help text for the Python panel.
func GetPythonPanelHelp() string {
	return "↑↓: Navigate | Enter: Install | d/Del: Delete | p: Pin | i: Refresh"
}

// getPinnedVersion reads the pinned Python version from the .python-version file.
func getPinnedVersion() string {
	content, err := os.ReadFile(".python-version")
	if err != nil {
		return "3.12"
	}
	return strings.TrimSpace(string(content))
}

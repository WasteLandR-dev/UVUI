// Package panels provides UI panels for the application.
package panels

import (
	"fmt"
	"strings"

	"uvui/internal/types"
	"uvui/internal/ui"
)

// ProjectState represents the project panel state.
type ProjectState struct {
	Name           string
	Status         *types.ProjectStatus
	Dependencies   []types.ProjectDependency
	DependencyTree *types.DependencyTree
	Selected       int
	Loading        bool
	ShowTree       bool
}

// RenderProjectPanel renders the project management panel.
func RenderProjectPanel(state *AppState) string {
	var content strings.Builder

	content.WriteString(ui.TitleStyle.Render("📦 Project Management"))
	content.WriteString("\n\n")

	if !state.Installed {
		content.WriteString(ui.ErrorStyle.Render("UV must be installed first to manage projects."))
		return content.String()
	}

	// Show loading state
	if state.ProjectState.Loading {
		content.WriteString(ui.LoadingStyle.Render("Loading project information..."))
		return content.String()
	}

	// Project status section
	content.WriteString(renderProjectStatus(state.ProjectState.Status))
	content.WriteString("\n")

	// Show project operations or initialization options
	if state.ProjectState.Status != nil && state.ProjectState.Status.IsProject {
		content.WriteString(renderProjectOperations(state))
		content.WriteString("\n")

		// Show dependencies or tree based on current view
		if state.ProjectState.ShowTree {
			content.WriteString(renderDependencyTree(state.ProjectState.DependencyTree))
		} else {
			content.WriteString(renderProjectDependencies(state.ProjectState.Dependencies))
		}
	} else {
		content.WriteString(renderInitializationHelp())
	}

	return content.String()
}

// renderProjectStatus renders the current project status.
func renderProjectStatus(status *types.ProjectStatus) string {
	var content strings.Builder

	content.WriteString(ui.CurrentVersionStyle.Render("Project Status"))
	content.WriteString("\n")

	if status == nil {
		content.WriteString(ui.ErrorStyle.Render("• Failed to detect project status"))
		return content.String()
	}

	if !status.IsProject {
		content.WriteString(ui.WarningMessageStyle.Render("• No project detected in current directory"))
		content.WriteString("\n")
		content.WriteString(ui.UnselectedItemStyle.Render(fmt.Sprintf("  Current directory: %s", status.Path)))
	} else {
		content.WriteString(ui.SuccessStyle.Render("• UV project detected"))
		content.WriteString("\n")
		content.WriteString(ui.InfoMessageStyle.Render(fmt.Sprintf("  Name: %s", status.Name)))
		content.WriteString("\n")
		content.WriteString(ui.UnselectedItemStyle.Render(fmt.Sprintf("  Path: %s", status.Path)))
		content.WriteString("\n")

		if status.PythonVersion != "" {
			content.WriteString(ui.InfoMessageStyle.Render(fmt.Sprintf("  Python: %s", status.PythonVersion)))
			content.WriteString("\n")
		}

		if status.HasLockFile {
			content.WriteString(ui.SuccessStyle.Render("  ✓ Dependencies locked (uv.lock found)"))
		} else {
			content.WriteString(ui.WarningMessageStyle.Render("  ⚠ Dependencies not locked"))
		}
		content.WriteString("\n")

		if status.HasVirtualEnv {
			content.WriteString(ui.SuccessStyle.Render("  ✓ Virtual environment exists"))
		} else {
			content.WriteString(ui.WarningMessageStyle.Render("  ⚠ No virtual environment found"))
		}
	}

	return content.String()
}

// renderProjectOperations renders available project operations.
func renderProjectOperations(_ *AppState) string {
	var content strings.Builder

	content.WriteString(ui.CurrentVersionStyle.Render("Available Operations"))
	content.WriteString("\n")

	operations := []struct {
		key         string
		description string
		available   bool
	}{
		{"s", "Sync dependencies", true},
		{"l", "Lock dependencies", true},
		{"t", "Toggle dependency tree view", true},
		{"r", "Refresh project status", true},
	}

	for _, op := range operations {
		style := ui.InfoMessageStyle
		if !op.available {
			style = ui.UnselectedItemStyle
		}
		content.WriteString(style.Render(fmt.Sprintf("  %s - %s", op.key, op.description)))
		content.WriteString("\n")
	}

	return content.String()
}

// renderProjectDependencies renders the project dependencies list.
func renderProjectDependencies(dependencies []types.ProjectDependency) string {
	var content strings.Builder

	content.WriteString(ui.CurrentVersionStyle.Render("Dependencies"))
	content.WriteString("\n")

	if len(dependencies) == 0 {
		content.WriteString(ui.UnselectedItemStyle.Render("  No dependencies found"))
		return content.String()
	}

	// Group dependencies by type
	mainDeps := []types.ProjectDependency{}
	devDeps := []types.ProjectDependency{}

	for _, dep := range dependencies {
		if dep.Type == "dev" {
			devDeps = append(devDeps, dep)
		} else {
			mainDeps = append(mainDeps, dep)
		}
	}

	// Render main dependencies
	if len(mainDeps) > 0 {
		content.WriteString(ui.InfoMessageStyle.Render("  Main:"))
		content.WriteString("\n")
		for _, dep := range mainDeps {
			content.WriteString(ui.UnselectedItemStyle.Render(fmt.Sprintf("    • %s %s", dep.Name, dep.Version)))
			content.WriteString("\n")
		}
	}

	// Render dev dependencies
	if len(devDeps) > 0 {
		content.WriteString(ui.InfoMessageStyle.Render("  Development:"))
		content.WriteString("\n")
		for _, dep := range devDeps {
			content.WriteString(ui.UnselectedItemStyle.Render(fmt.Sprintf("    • %s %s", dep.Name, dep.Version)))
			content.WriteString("\n")
		}
	}

	return content.String()
}

// renderDependencyTree renders the dependency tree.
func renderDependencyTree(tree *types.DependencyTree) string {
	var content strings.Builder

	content.WriteString(ui.CurrentVersionStyle.Render("Dependency Tree"))
	content.WriteString("\n")

	if tree == nil || len(tree.Dependencies) == 0 {
		content.WriteString(ui.UnselectedItemStyle.Render("  No dependency tree available"))
		return content.String()
	}

	for _, node := range tree.Dependencies {
		indent := strings.Repeat("  ", node.Level)
		prefix := "├─"
		if node.Level == 0 {
			prefix = ""
		}

		line := fmt.Sprintf("%s%s %s", indent, prefix, node.Name)
		if node.Version != "" {
			line += fmt.Sprintf(" (%s)", node.Version)
		}

		// Use different styling for different levels
		if node.Level == 0 {
			content.WriteString(ui.InfoMessageStyle.Render(line))
		} else {
			content.WriteString(ui.UnselectedItemStyle.Render(line))
		}
		content.WriteString("\n")
	}

	return content.String()
}

// renderInitializationHelp renders help for project initialization.
func renderInitializationHelp() string {
	var content strings.Builder

	content.WriteString(ui.CurrentVersionStyle.Render("Initialize New Project"))
	content.WriteString("\n")

	content.WriteString(ui.InfoMessageStyle.Render("  Available commands:"))
	content.WriteString("\n")
	content.WriteString(ui.UnselectedItemStyle.Render("    i - Initialize new project in current directory"))
	content.WriteString("\n")
	content.WriteString(ui.UnselectedItemStyle.Render("    n - Initialize new project with custom name"))
	content.WriteString("\n")
	content.WriteString(ui.UnselectedItemStyle.Render("    a - Initialize as application project"))
	content.WriteString("\n")
	content.WriteString(ui.UnselectedItemStyle.Render("    l - Initialize as library project"))
	content.WriteString("\n")
	content.WriteString(ui.UnselectedItemStyle.Render("    r - Refresh project status"))
	content.WriteString("\n\n")

	content.WriteString(ui.WarningMessageStyle.Render("Note: Project initialization will create pyproject.toml and basic structure"))

	return content.String()
}

// GetProjectPanelHelp returns help text for the project panel.
func GetProjectPanelHelp() string {
	help := []string{
		"Project Management Panel Help:",
		"",
		"When no project detected:",
		"  i - Initialize new project",
		"  n - Initialize with custom name",
		"  a - Initialize as app",
		"  l - Initialize as library",
		"",
		"When in project:",
		"  s - Sync dependencies",
		"  l - Lock dependencies",
		"  t - Toggle tree view",
		"  r - Refresh status",
		"",
		"Navigation:",
		"  ↑/↓ - Navigate items",
		"  Tab - Switch panels",
		"  ? - Show/hide help",
		"  q - Quit",
	}

	return strings.Join(help, "\n")
}

// Package panels provides UI panels for the application.
package panels

import (
	"uvui/internal/ui"
)

// RenderEnvironmentPanel renders the environment management panel
func RenderEnvironmentPanel(state *AppState) string {
	content := "Environment Management\n\n"

	if !state.Installed {
		content += ui.ErrorStyle.Render("UV must be installed first to manage environments.")
		return content
	}

	content += "Virtual environments:\n"
	content += "• No virtual environments found\n"
	content += "• Use 'uv venv' to create environments (not implemented yet)\n"

	return content
}

// GetEnvironmentPanelHelp returns help text for the environment panel
func GetEnvironmentPanelHelp() string {
	return "Coming in Phase 5: Environment operations"
}

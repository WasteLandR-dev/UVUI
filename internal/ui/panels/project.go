package panels

import (
	"uvui/internal/types"
	"uvui/internal/ui"
)

// RenderProjectPanel renders the project management panel
func RenderProjectPanel(state *types.AppState) string {
	content := "Project Management\n\n"

	if !state.UVStatus.Installed {
		content += ui.ErrorStyle.Render("UV must be installed first to manage projects.")
		return content
	}

	content += "Current project status:\n"
	content += "• No project detected in current directory\n"
	content += "• Use 'uv init' to create a new project (not implemented yet)\n"

	return content
}

// GetProjectPanelHelp returns help text for the project panel
func GetProjectPanelHelp() string {
	return "Coming in Phase 3: Project operations"
}

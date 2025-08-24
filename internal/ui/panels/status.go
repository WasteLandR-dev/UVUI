package panels

import (
	"fmt"
	"strings"

	"uvui/internal/services"
	"uvui/internal/types"
	"uvui/internal/ui"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RenderStatusPanel renders the status panel content
func RenderStatusPanel(state *types.AppState, uvInstaller services.UVInstallerInterface) string {
	var content strings.Builder

	content.WriteString("UV Installation Status:\n\n")

	if state.Installing {
		content.WriteString(ui.LoadingStyle.Render("⏳ Installing UV...") + "\n")
	} else if state.UVStatus.Installed {
		content.WriteString(ui.SuccessStyle.Render("✅ UV is installed") + "\n")
		if state.UVStatus.Version != "" {
			content.WriteString(fmt.Sprintf("   Version: %s\n", state.UVStatus.Version))
		}
		if state.UVStatus.Path != "" {
			content.WriteString(fmt.Sprintf("   Path: %s\n", state.UVStatus.Path))
		}
	} else {
		content.WriteString(ui.ErrorStyle.Render("❌ UV is not installed") + "\n")
		content.WriteString("   Press 'i' to install UV\n\n")

		// Show install command
		cmd, err := uvInstaller.GetInstallCommand()
		if err == nil {
			content.WriteString("Install command:\n")
			content.WriteString(fmt.Sprintf("   %s\n", cmd))
		}
	}

	// Show recent messages
	if len(state.Messages) > 0 {
		content.WriteString("\nRecent Messages:\n")
		for _, msg := range state.Messages {
			content.WriteString(fmt.Sprintf("• %s\n", msg))
		}
	}

	// Show current operation status
	if state.Operation.InProgress {
		content.WriteString("\nCurrent Operation:\n")
		content.WriteString(ui.LoadingStyle.Render(fmt.Sprintf("⏳ %s %s...",
			cases.Title(language.English).String(state.Operation.Operation),
			state.Operation.Target)))
	}

	return content.String()
}

// GetStatusPanelHelp returns help text for the status panel
func GetStatusPanelHelp() string {
	return "i: Install UV | r: Refresh status | ?: Help"
}

package panels

import (
	"errors"
	"strings"
	"testing"

	"uvui/internal/types"
)

func TestRenderStatusPanel_Installing(t *testing.T) {
	state := &AppState{
		Installing: true,
		UVStatus: types.UVStatus{
			Installed: false,
		},
	}

	result := RenderStatusPanel(state, "")

	// Check that it contains the installation status header
	if !strings.Contains(result, "UV Installation Status:") {
		t.Error("Expected status header not found")
	}

	// Check that it shows installing message
	expected := "⏳ Installing UV..."
	if !strings.Contains(result, expected) {
		t.Errorf("Expected installing message '%s' not found in result", expected)
	}

	// Should not contain installed or not installed messages
	if strings.Contains(result, "✅ UV is installed") {
		t.Error("Should not show installed message when installing")
	}
	if strings.Contains(result, "❌ UV is not installed") {
		t.Error("Should not show not installed message when installing")
	}
}

func TestRenderStatusPanel_Installed_WithVersionAndPath(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
			Version:   "1.2.3",
			Path:      "/usr/local/bin/uv",
		},
	}

	result := RenderStatusPanel(state, "")

	// Check installed message
	if !strings.Contains(result, "✅ UV is installed") {
		t.Error("Expected installed message not found")
	}

	// Check version is displayed
	if !strings.Contains(result, "Version: 1.2.3") {
		t.Error("Expected version not found")
	}

	// Check path is displayed
	if !strings.Contains(result, "Path: /usr/local/bin/uv") {
		t.Error("Expected path not found")
	}

	// Should not contain installing or not installed messages
	if strings.Contains(result, "⏳ Installing UV...") {
		t.Error("Should not show installing message when installed")
	}
	if strings.Contains(result, "❌ UV is not installed") {
		t.Error("Should not show not installed message when installed")
	}
}

func TestRenderStatusPanel_Installed_WithoutVersionAndPath(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
			Version:   "",
			Path:      "",
		},
	}

	result := RenderStatusPanel(state, "")

	// Check installed message
	if !strings.Contains(result, "✅ UV is installed") {
		t.Error("Expected installed message not found")
	}

	// Should not display version or path when empty
	if strings.Contains(result, "Version:") {
		t.Error("Should not show version when empty")
	}
	if strings.Contains(result, "Path:") {
		t.Error("Should not show path when empty")
	}
}

func TestRenderStatusPanel_NotInstalled_WithInstallCommand(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: false,
		},
	}
	installCmd := "curl -LsSf https://astral.sh/uv/install.sh | sh"

	result := RenderStatusPanel(state, installCmd)

	// Check not installed message
	if !strings.Contains(result, "❌ UV is not installed") {
		t.Error("Expected not installed message not found")
	}

	// Check install instruction
	if !strings.Contains(result, "Press 'i' to install UV") {
		t.Error("Expected install instruction not found")
	}

	// Check install command is displayed
	if !strings.Contains(result, "Install command:") {
		t.Error("Expected install command header not found")
	}
	if !strings.Contains(result, installCmd) {
		t.Error("Expected install command not found")
	}
}

func TestRenderStatusPanel_NotInstalled_WithoutInstallCommand(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: false,
		},
	}

	result := RenderStatusPanel(state, "")

	// Check not installed message
	if !strings.Contains(result, "❌ UV is not installed") {
		t.Error("Expected not installed message not found")
	}

	// Check install instruction
	if !strings.Contains(result, "Press 'i' to install UV") {
		t.Error("Expected install instruction not found")
	}

	// Should not show install command section when empty
	if strings.Contains(result, "Install command:") {
		t.Error("Should not show install command header when command is empty")
	}
}

func TestRenderStatusPanel_WithMessages(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
		},
		Messages: []string{"Package installed successfully", "Configuration updated", "Ready to use"},
	}

	result := RenderStatusPanel(state, "")

	// Check messages header
	if !strings.Contains(result, "Recent Messages:") {
		t.Error("Expected messages header not found")
	}

	// Check all messages are displayed
	for _, msg := range state.Messages {
		expected := "• " + msg
		if !strings.Contains(result, expected) {
			t.Errorf("Expected message '%s' not found", expected)
		}
	}
}

func TestRenderStatusPanel_WithoutMessages(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
		},
		Messages: []string{},
	}

	result := RenderStatusPanel(state, "")

	// Should not show messages section when empty
	if strings.Contains(result, "Recent Messages:") {
		t.Error("Should not show messages header when no messages")
	}
}

func TestRenderStatusPanel_WithCurrentOperation(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
		},
		Operation: types.OperationStatus{
			InProgress: true,
			Operation:  "installing",
			Target:     "package-name",
		},
	}

	result := RenderStatusPanel(state, "")

	// Check operation header
	if !strings.Contains(result, "Current Operation:") {
		t.Error("Expected operation header not found")
	}

	// Check operation details (should be title-cased)
	expected := "⏳ Installing package-name..."
	if !strings.Contains(result, expected) {
		t.Errorf("Expected operation message '%s' not found", expected)
	}
}

func TestRenderStatusPanel_WithoutCurrentOperation(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
		},
		Operation: types.OperationStatus{
			InProgress: false,
		},
	}

	result := RenderStatusPanel(state, "")

	// Should not show operation section when not in progress
	if strings.Contains(result, "Current Operation:") {
		t.Error("Should not show operation header when no operation in progress")
	}
}

func TestRenderStatusPanel_ComplexScenario(t *testing.T) {
	// Test with all possible elements present
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
			Version:   "2.0.1",
			Path:      "/opt/uv/bin/uv",
		},
		Messages: []string{"Started successfully", "All systems operational"},
		Operation: types.OperationStatus{
			InProgress: true,
			Operation:  "updating",
			Target:     "dependencies",
		},
	}
	installCmd := "brew install uv"

	result := RenderStatusPanel(state, installCmd)

	// Verify all sections are present
	sections := []string{
		"UV Installation Status:",
		"✅ UV is installed",
		"Version: 2.0.1",
		"Path: /opt/uv/bin/uv",
		"Recent Messages:",
		"• Started successfully",
		"• All systems operational",
		"Current Operation:",
		"⏳ Updating dependencies...",
	}

	for _, section := range sections {
		if !strings.Contains(result, section) {
			t.Errorf("Expected section '%s' not found in complex scenario", section)
		}
	}
}

func TestRenderStatusPanel_OperationTitleCase(t *testing.T) {
	// Test that operation names are properly title-cased
	testCases := []struct {
		operation string
		expected  string
	}{
		{"installing", "Installing"},
		{"updating", "Updating"},
		{"removing", "Removing"},
		{"configuring", "Configuring"},
	}

	for _, tc := range testCases {
		state := &AppState{
			Operation: types.OperationStatus{
				InProgress: true,
				Operation:  tc.operation,
				Target:     "test-target",
			},
		}

		result := RenderStatusPanel(state, "")
		expected := "⏳ " + tc.expected + " test-target..."

		if !strings.Contains(result, expected) {
			t.Errorf("For operation '%s', expected '%s' but not found in result", tc.operation, expected)
		}
	}
}

func TestRenderStatusPanel_OperationWithAllFields(t *testing.T) {
	// Test operation with Success and Error fields (even though they're not used in render)
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
		},
		Operation: types.OperationStatus{
			InProgress: true,
			Operation:  "testing",
			Target:     "comprehensive",
			Success:    true,
			Error:      errors.New("test error"),
		},
	}

	result := RenderStatusPanel(state, "")

	// Should still show operation correctly regardless of Success/Error fields
	expected := "⏳ Testing comprehensive..."
	if !strings.Contains(result, expected) {
		t.Error("Operation should be displayed correctly regardless of Success/Error fields")
	}
}

func TestGetStatusPanelHelp(t *testing.T) {
	result := GetStatusPanelHelp()
	expected := "i: Install UV | r: Refresh status | ?: Help"

	if result != expected {
		t.Errorf("Expected help text '%s', got '%s'", expected, result)
	}
}

// Test edge cases and boundary conditions

func TestRenderStatusPanel_EmptyState(t *testing.T) {
	state := &AppState{}
	result := RenderStatusPanel(state, "")

	// Should show not installed state by default (UVStatus.Installed defaults to false)
	if !strings.Contains(result, "❌ UV is not installed") {
		t.Error("Empty state should show not installed message")
	}

	if !strings.Contains(result, "UV Installation Status:") {
		t.Error("Should always show status header")
	}
}

func TestRenderStatusPanel_NilMessages(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
		},
		Messages: nil, // explicitly nil
	}

	result := RenderStatusPanel(state, "")

	// Should not crash and should not show messages section
	if strings.Contains(result, "Recent Messages:") {
		t.Error("Should not show messages header when messages is nil")
	}
}

func TestRenderStatusPanel_EmptyStrings(t *testing.T) {
	state := &AppState{
		Installing: false,
		UVStatus: types.UVStatus{
			Installed: true,
			Version:   "",
			Path:      "",
		},
		Messages: []string{"", "  ", "valid message"},
	}

	result := RenderStatusPanel(state, "")

	// Should handle empty strings gracefully
	if strings.Contains(result, "Version:") {
		t.Error("Should not show version when empty string")
	}
	if strings.Contains(result, "Path:") {
		t.Error("Should not show path when empty string")
	}

	// Should still show messages section and include all messages (even empty ones)
	if !strings.Contains(result, "Recent Messages:") {
		t.Error("Should show messages header when messages slice is not empty")
	}
	if !strings.Contains(result, "• valid message") {
		t.Error("Should show valid message")
	}
}

func TestRenderStatusPanel_AllUVStatusFields(t *testing.T) {
	// Test all combinations of UVStatus fields
	testCases := []struct {
		name       string
		uvStatus   types.UVStatus
		shouldShow map[string]bool
	}{
		{
			name: "All fields populated",
			uvStatus: types.UVStatus{
				Installed: true,
				Version:   "1.0.0",
				Path:      "/usr/bin/uv",
			},
			shouldShow: map[string]bool{
				"✅ UV is installed": true,
				"Version: 1.0.0":    true,
				"Path: /usr/bin/uv": true,
			},
		},
		{
			name: "Only installed",
			uvStatus: types.UVStatus{
				Installed: true,
			},
			shouldShow: map[string]bool{
				"✅ UV is installed": true,
				"Version:":          false,
				"Path:":             false,
			},
		},
		{
			name: "Not installed",
			uvStatus: types.UVStatus{
				Installed: false,
				Version:   "should-not-show",
				Path:      "should-not-show",
			},
			shouldShow: map[string]bool{
				"❌ UV is not installed": true,
				"Version:":              false,
				"Path:":                 false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := &AppState{
				UVStatus: tc.uvStatus,
			}

			result := RenderStatusPanel(state, "")

			for text, shouldBePresent := range tc.shouldShow {
				if shouldBePresent && !strings.Contains(result, text) {
					t.Errorf("Expected '%s' to be present but it wasn't", text)
				}
				if !shouldBePresent && strings.Contains(result, text) {
					t.Errorf("Expected '%s' to not be present but it was", text)
				}
			}
		})
	}
}

func TestRenderStatusPanel_BooleanLogicCoverage(t *testing.T) {
	// Test specific boolean combinations to ensure all branches are covered
	testCases := []struct {
		name        string
		installing  bool
		installed   bool
		expectedMsg string
	}{
		{
			name:        "Installing true, installed false",
			installing:  true,
			installed:   false,
			expectedMsg: "⏳ Installing UV...",
		},
		{
			name:        "Installing false, installed true",
			installing:  false,
			installed:   true,
			expectedMsg: "✅ UV is installed",
		},
		{
			name:        "Installing false, installed false",
			installing:  false,
			installed:   false,
			expectedMsg: "❌ UV is not installed",
		},
		{
			name:        "Installing true, installed true (installing takes precedence)",
			installing:  true,
			installed:   true,
			expectedMsg: "⏳ Installing UV...",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := &AppState{
				Installing: tc.installing,
				UVStatus: types.UVStatus{
					Installed: tc.installed,
				},
			}

			result := RenderStatusPanel(state, "")

			if !strings.Contains(result, tc.expectedMsg) {
				t.Errorf("Expected message '%s' not found in result for case '%s'", tc.expectedMsg, tc.name)
			}
		})
	}
}

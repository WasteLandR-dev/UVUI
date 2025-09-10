package panels

import (
	"strings"
	"testing"

	"uvui/internal/types"
)

func TestRenderProjectPanel_UVNotInstalled(t *testing.T) {
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: false,
		},
	}

	result := RenderProjectPanel(state)

	// Should show title
	if !strings.Contains(result, "📦 Project Management") {
		t.Error("Expected project management title not found")
	}

	// Should show UV not installed error
	if !strings.Contains(result, "UV must be installed first to manage projects.") {
		t.Error("Expected UV installation error message not found")
	}

	// Should not show any other content
	if strings.Contains(result, "Project Status") {
		t.Error("Should not show project status when UV not installed")
	}
}

func TestRenderProjectPanel_Loading(t *testing.T) {
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Loading: true,
		},
	}

	result := RenderProjectPanel(state)

	// Should show loading message
	if !strings.Contains(result, "Loading project information...") {
		t.Error("Expected loading message not found")
	}

	// Should not show other content
	if strings.Contains(result, "Project Status") {
		t.Error("Should not show project status when loading")
	}
}

func TestRenderProjectPanel_ProjectDetected_WithTree(t *testing.T) {
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Loading:  false,
			ShowTree: true,
			Status: &types.ProjectStatus{
				IsProject: true,
				Name:      "test-project",
				Path:      "/path/to/project",
			},
			DependencyTree: &types.DependencyTree{
				Dependencies: []types.TreeNode{
					{Name: "requests", Version: "2.28.1", Level: 0},
				},
			},
		},
	}

	result := RenderProjectPanel(state)

	// Should show project status
	if !strings.Contains(result, "Project Status") {
		t.Error("Expected project status section not found")
	}

	// Should show available operations
	if !strings.Contains(result, "Available Operations") {
		t.Error("Expected available operations section not found")
	}

	// Should show dependency tree (not dependencies list)
	if !strings.Contains(result, "Dependency Tree") {
		t.Error("Expected dependency tree section not found")
	}

	// Should not show dependencies list
	if strings.Contains(result, "Dependencies") && !strings.Contains(result, "Dependency Tree") {
		t.Error("Should not show dependencies list when showing tree")
	}
}

func TestRenderProjectPanel_ProjectDetected_WithDependencies(t *testing.T) {
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Loading:  false,
			ShowTree: false,
			Status: &types.ProjectStatus{
				IsProject: true,
				Name:      "test-project",
			},
			Dependencies: []types.ProjectDependency{
				{Name: "requests", Version: "2.28.1", Type: "main"},
			},
		},
	}

	result := RenderProjectPanel(state)

	// Should show project status
	if !strings.Contains(result, "Project Status") {
		t.Error("Expected project status section not found")
	}

	// Should show available operations
	if !strings.Contains(result, "Available Operations") {
		t.Error("Expected available operations section not found")
	}

	// Should show dependencies (not tree)
	if !strings.Contains(result, "Dependencies") {
		t.Error("Expected dependencies section not found")
	}

	// Should not show dependency tree
	if strings.Contains(result, "Dependency Tree") {
		t.Error("Should not show dependency tree when showing dependencies list")
	}
}

func TestRenderProjectPanel_NoProjectDetected(t *testing.T) {
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Loading: false,
			Status: &types.ProjectStatus{
				IsProject: false,
				Path:      "/current/path",
			},
		},
	}

	result := RenderProjectPanel(state)

	// Should show project status
	if !strings.Contains(result, "Project Status") {
		t.Error("Expected project status section not found")
	}

	// Should show initialization help instead of operations
	if !strings.Contains(result, "Initialize New Project") {
		t.Error("Expected initialization help section not found")
	}

	// Should not show available operations
	if strings.Contains(result, "Available Operations") {
		t.Error("Should not show operations when no project detected")
	}
}

func TestRenderProjectStatus_NilStatus(t *testing.T) {
	result := renderProjectStatus(nil)

	if !strings.Contains(result, "Project Status") {
		t.Error("Expected project status header not found")
	}

	if !strings.Contains(result, "Failed to detect project status") {
		t.Error("Expected error message for nil status not found")
	}
}

func TestRenderProjectStatus_NoProject(t *testing.T) {
	status := &types.ProjectStatus{
		IsProject: false,
		Path:      "/current/directory",
	}

	result := renderProjectStatus(status)

	if !strings.Contains(result, "No project detected in current directory") {
		t.Error("Expected no project message not found")
	}

	if !strings.Contains(result, "Current directory: /current/directory") {
		t.Error("Expected current directory path not found")
	}
}

func TestRenderProjectStatus_ProjectDetected_FullInfo(t *testing.T) {
	status := &types.ProjectStatus{
		IsProject:     true,
		Name:          "my-awesome-project",
		Path:          "/path/to/project",
		PythonVersion: "3.11.0",
		HasLockFile:   true,
		HasVirtualEnv: true,
	}

	result := renderProjectStatus(status)

	expectedStrings := []string{
		"UV project detected",
		"Name: my-awesome-project",
		"Path: /path/to/project",
		"Python: 3.11.0",
		"Dependencies locked (uv.lock found)",
		"Virtual environment exists",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected string '%s' not found in result", expected)
		}
	}
}

func TestRenderProjectStatus_ProjectDetected_MinimalInfo(t *testing.T) {
	status := &types.ProjectStatus{
		IsProject:     true,
		Name:          "minimal-project",
		Path:          "/minimal/path",
		PythonVersion: "", // Empty python version
		HasLockFile:   false,
		HasVirtualEnv: false,
	}

	result := renderProjectStatus(status)

	// Should show basic info
	if !strings.Contains(result, "Name: minimal-project") {
		t.Error("Expected project name not found")
	}
	if !strings.Contains(result, "Path: /minimal/path") {
		t.Error("Expected project path not found")
	}

	// Should not show Python version when empty
	if strings.Contains(result, "Python:") {
		t.Error("Should not show Python version when empty")
	}

	// Should show warnings for missing lock file and virtual env
	if !strings.Contains(result, "Dependencies not locked") {
		t.Error("Expected dependencies not locked warning not found")
	}
	if !strings.Contains(result, "No virtual environment found") {
		t.Error("Expected no virtual environment warning not found")
	}
}

func TestRenderProjectOperations(t *testing.T) {
	state := &AppState{}
	result := renderProjectOperations(state)

	if !strings.Contains(result, "Available Operations") {
		t.Error("Expected operations header not found")
	}

	// Check all operations are listed
	operations := []string{
		"s - Sync dependencies",
		"l - Lock dependencies",
		"t - Toggle dependency tree view",
		"r - Refresh project status",
	}

	for _, op := range operations {
		if !strings.Contains(result, op) {
			t.Errorf("Expected operation '%s' not found", op)
		}
	}
}

func TestRenderProjectDependencies_Empty(t *testing.T) {
	dependencies := []types.ProjectDependency{}
	result := renderProjectDependencies(dependencies)

	if !strings.Contains(result, "Dependencies") {
		t.Error("Expected dependencies header not found")
	}

	if !strings.Contains(result, "No dependencies found") {
		t.Error("Expected no dependencies message not found")
	}
}

func TestRenderProjectDependencies_MainOnly(t *testing.T) {
	dependencies := []types.ProjectDependency{
		{Name: "requests", Version: "2.28.1", Type: "main"},
		{Name: "click", Version: "8.1.3", Type: "main"},
	}

	result := renderProjectDependencies(dependencies)

	if !strings.Contains(result, "Main:") {
		t.Error("Expected main dependencies section not found")
	}

	if !strings.Contains(result, "• requests 2.28.1") {
		t.Error("Expected requests dependency not found")
	}

	if !strings.Contains(result, "• click 8.1.3") {
		t.Error("Expected click dependency not found")
	}

	// Should not show dev section
	if strings.Contains(result, "Development:") {
		t.Error("Should not show development section when no dev dependencies")
	}
}

func TestRenderProjectDependencies_DevOnly(t *testing.T) {
	dependencies := []types.ProjectDependency{
		{Name: "pytest", Version: "7.2.0", Type: "dev"},
		{Name: "black", Version: "22.10.0", Type: "dev"},
	}

	result := renderProjectDependencies(dependencies)

	if !strings.Contains(result, "Development:") {
		t.Error("Expected development dependencies section not found")
	}

	if !strings.Contains(result, "• pytest 7.2.0") {
		t.Error("Expected pytest dependency not found")
	}

	if !strings.Contains(result, "• black 22.10.0") {
		t.Error("Expected black dependency not found")
	}

	// Should not show main section
	if strings.Contains(result, "Main:") {
		t.Error("Should not show main section when no main dependencies")
	}
}

func TestRenderProjectDependencies_Mixed(t *testing.T) {
	dependencies := []types.ProjectDependency{
		{Name: "requests", Version: "2.28.1", Type: "main"},
		{Name: "pytest", Version: "7.2.0", Type: "dev"},
		{Name: "click", Version: "8.1.3", Type: ""}, // Should be treated as main
		{Name: "black", Version: "22.10.0", Type: "dev"},
	}

	result := renderProjectDependencies(dependencies)

	// Should show both sections
	if !strings.Contains(result, "Main:") {
		t.Error("Expected main dependencies section not found")
	}
	if !strings.Contains(result, "Development:") {
		t.Error("Expected development dependencies section not found")
	}

	// Check main dependencies (including empty type)
	if !strings.Contains(result, "• requests 2.28.1") {
		t.Error("Expected requests in main dependencies")
	}
	if !strings.Contains(result, "• click 8.1.3") {
		t.Error("Expected click in main dependencies (empty type)")
	}

	// Check dev dependencies
	if !strings.Contains(result, "• pytest 7.2.0") {
		t.Error("Expected pytest in dev dependencies")
	}
	if !strings.Contains(result, "• black 22.10.0") {
		t.Error("Expected black in dev dependencies")
	}
}

func TestRenderDependencyTree_Nil(t *testing.T) {
	result := renderDependencyTree(nil)

	if !strings.Contains(result, "Dependency Tree") {
		t.Error("Expected dependency tree header not found")
	}

	if !strings.Contains(result, "No dependency tree available") {
		t.Error("Expected no tree message for nil tree not found")
	}
}

func TestRenderDependencyTree_Empty(t *testing.T) {
	tree := &types.DependencyTree{
		Dependencies: []types.TreeNode{},
	}

	result := renderDependencyTree(tree)

	if !strings.Contains(result, "Dependency Tree") {
		t.Error("Expected dependency tree header not found")
	}

	if !strings.Contains(result, "No dependency tree available") {
		t.Error("Expected no tree message for empty tree not found")
	}
}

func TestRenderDependencyTree_WithNodes(t *testing.T) {
	tree := &types.DependencyTree{
		Dependencies: []types.TreeNode{
			{Name: "requests", Version: "2.28.1", Level: 0},
			{Name: "urllib3", Version: "1.26.12", Level: 1},
			{Name: "certifi", Version: "2022.9.24", Level: 1},
			{Name: "charset-normalizer", Version: "2.1.1", Level: 1},
		},
	}

	result := renderDependencyTree(tree)

	if !strings.Contains(result, "Dependency Tree") {
		t.Error("Expected dependency tree header not found")
	}

	// Check root level (no prefix, no indent)
	if !strings.Contains(result, "requests (2.28.1)") {
		t.Error("Expected root dependency not found")
	}

	// Check level 1 dependencies (with prefix and indent)
	expectedLevel1 := []string{
		"├─ urllib3 (1.26.12)",
		"├─ certifi (2022.9.24)",
		"├─ charset-normalizer (2.1.1)",
	}

	for _, expected := range expectedLevel1 {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected level 1 dependency '%s' not found", expected)
		}
	}
}

func TestRenderDependencyTree_WithoutVersions(t *testing.T) {
	tree := &types.DependencyTree{
		Dependencies: []types.TreeNode{
			{Name: "root-package", Version: "", Level: 0},
			{Name: "child-package", Version: "", Level: 1},
		},
	}

	result := renderDependencyTree(tree)

	// Should show names without version info
	if !strings.Contains(result, "root-package\n") {
		t.Error("Expected root package without version not found")
	}

	if !strings.Contains(result, "├─ child-package") {
		t.Error("Expected child package without version not found")
	}

	// Should not contain version parentheses
	if strings.Contains(result, "()") {
		t.Error("Should not show empty version parentheses")
	}
}

func TestRenderDependencyTree_MultipleRoots(t *testing.T) {
	tree := &types.DependencyTree{
		Dependencies: []types.TreeNode{
			{Name: "requests", Version: "2.28.1", Level: 0},
			{Name: "urllib3", Version: "1.26.12", Level: 1},
			{Name: "flask", Version: "2.2.2", Level: 0},
			{Name: "jinja2", Version: "3.1.2", Level: 1},
		},
	}

	result := renderDependencyTree(tree)

	// Should show both root packages
	if !strings.Contains(result, "requests (2.28.1)") {
		t.Error("Expected first root package not found")
	}
	if !strings.Contains(result, "flask (2.2.2)") {
		t.Error("Expected second root package not found")
	}

	// Should show children with proper indentation
	if !strings.Contains(result, "├─ urllib3 (1.26.12)") {
		t.Error("Expected first child dependency not found")
	}
	if !strings.Contains(result, "├─ jinja2 (3.1.2)") {
		t.Error("Expected second child dependency not found")
	}
}

func TestRenderInitializationHelp(t *testing.T) {
	result := renderInitializationHelp()

	if !strings.Contains(result, "Initialize New Project") {
		t.Error("Expected initialization header not found")
	}

	// Check all initialization commands
	commands := []string{
		"i - Initialize new project in current directory",
		"n - Initialize new project with custom name",
		"a - Initialize as application project",
		"l - Initialize as library project",
		"r - Refresh project status",
	}

	for _, cmd := range commands {
		if !strings.Contains(result, cmd) {
			t.Errorf("Expected command '%s' not found", cmd)
		}
	}

	if !strings.Contains(result, "Note: Project initialization will create pyproject.toml and basic structure") {
		t.Error("Expected initialization note not found")
	}
}

func TestGetProjectPanelHelp(t *testing.T) {
	result := GetProjectPanelHelp()

	// Check main sections
	sections := []string{
		"Project Management Panel Help:",
		"When no project detected:",
		"When in project:",
		"Navigation:",
	}

	for _, section := range sections {
		if !strings.Contains(result, section) {
			t.Errorf("Expected help section '%s' not found", section)
		}
	}

	// Check specific commands
	commands := []string{
		"i - Initialize new project",
		"n - Initialize with custom name",
		"s - Sync dependencies",
		"l - Lock dependencies",
		"t - Toggle tree view",
		"q - Quit",
	}

	for _, cmd := range commands {
		if !strings.Contains(result, cmd) {
			t.Errorf("Expected help command '%s' not found", cmd)
		}
	}
}

// Edge cases and complex scenarios

func TestRenderProjectPanel_ComplexScenario(t *testing.T) {
	// Test a complex scenario with all features enabled
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Loading:  false,
			ShowTree: false,
			Status: &types.ProjectStatus{
				IsProject:     true,
				Name:          "complex-project",
				Path:          "/path/to/complex",
				PythonVersion: "3.11.5",
				HasLockFile:   true,
				HasVirtualEnv: true,
			},
			Dependencies: []types.ProjectDependency{
				{Name: "requests", Version: "2.28.1", Type: "main"},
				{Name: "pytest", Version: "7.2.0", Type: "dev"},
			},
			DependencyTree: &types.DependencyTree{
				Dependencies: []types.TreeNode{
					{Name: "requests", Version: "2.28.1", Level: 0},
				},
			},
		},
	}

	result := RenderProjectPanel(state)

	// Should contain all major sections
	majorSections := []string{
		"📦 Project Management",
		"Project Status",
		"UV project detected",
		"Name: complex-project",
		"Available Operations",
		"Dependencies",
		"Main:",
		"Development:",
	}

	for _, section := range majorSections {
		if !strings.Contains(result, section) {
			t.Errorf("Expected section '%s' not found in complex scenario", section)
		}
	}
}

func TestRenderProjectPanel_EdgeCases(t *testing.T) {
	// Test with nil project status
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Loading: false,
			Status:  nil, // This should trigger the nil status handling
		},
	}

	result := RenderProjectPanel(state)

	if !strings.Contains(result, "Failed to detect project status") {
		t.Error("Expected error message for nil project status")
	}

	// Should show initialization help when status is nil
	if !strings.Contains(result, "Initialize New Project") {
		t.Error("Expected initialization help when status is nil")
	}
}

func TestProjectState_AllFields(t *testing.T) {
	// Test that all ProjectState fields are properly handled
	state := &AppState{
		UVStatus: types.UVStatus{
			Installed: true,
		},
		ProjectState: ProjectState{
			Name:     "test-name", // This field isn't used in current implementation
			Selected: 5,           // This field isn't used in current implementation
			Loading:  false,
			ShowTree: true,
			Status: &types.ProjectStatus{
				IsProject: true,
				Name:      "actual-project-name",
			},
			Dependencies:   []types.ProjectDependency{},
			DependencyTree: &types.DependencyTree{Dependencies: []types.TreeNode{}},
		},
	}

	result := RenderProjectPanel(state)

	// Even though Name and Selected aren't used, the function should still work
	if !strings.Contains(result, "📦 Project Management") {
		t.Error("Should render properly even with unused ProjectState fields")
	}

	// Should show dependency tree view since ShowTree is true
	if !strings.Contains(result, "Dependency Tree") {
		t.Error("Should show dependency tree when ShowTree is true")
	}
}

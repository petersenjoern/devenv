package tui

import (
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
)

func TestShowEnvironmentSelection_ShouldReturnValidEnvironment(t *testing.T) {
	cfg := config.Config{}
	tui := New(cfg)
	
	// This test would normally require user interaction
	// For now, test that it returns a valid environment option
	env, err := tui.ShowEnvironmentSelection()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	validEnvs := []string{"linux", "wsl"}
	isValid := false
	for _, validEnv := range validEnvs {
		if env == validEnv {
			isValid = true
			break
		}
	}
	
	if !isValid {
		t.Errorf("Expected environment to be 'linux' or 'wsl', got: %s", env)
	}
}

func TestShowToolSelection_ShouldReturnSelectedTools(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
				"zsh": config.ToolConfig{
					DisplayName: "Zsh Shell", 
					BinaryName:  "zsh",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	selectedTools, err := tui.ShowToolSelection()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Should return a slice of selected tool names
	if selectedTools == nil {
		t.Errorf("Expected selectedTools to not be nil")
	}
}

func TestShowToolSelectionByCategory_ShouldOrganizeToolsByCategory(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
				"zsh": config.ToolConfig{
					DisplayName: "Zsh Shell", 
					BinaryName:  "zsh",
				},
			},
			"editors": {
				"vim": config.ToolConfig{
					DisplayName: "Vim Editor",
					BinaryName:  "vim",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	selections, err := tui.ShowToolSelectionByCategory()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if selections.CategoryAndTools == nil {
		t.Errorf("Expected selections.CategoryAndTools to not be nil")
	}
	
	// Should have entries for both categories
	if len(selections.CategoryAndTools) == 0 {
		t.Errorf("Expected selections to contain category entries, got empty")
	}
	
	// Should contain at least the categories we defined
	foundShells := false
	foundEditors := false
	for _, categoryTools := range selections.CategoryAndTools {
		if categoryTools.Category == "shells" {
			foundShells = true
		}
		if categoryTools.Category == "editors" {
			foundEditors = true
		}
	}
	
	if !foundShells {
		t.Errorf("Expected to find 'shells' category in selections")
	}
	if !foundEditors {
		t.Errorf("Expected to find 'editors' category in selections")
	}
}

func TestDetectActualEnvironment_ShouldDetectWSLOrLinux(t *testing.T) {
	cfg := config.Config{}
	tui := New(cfg)
	
	env, err := tui.DetectActualEnvironment()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Should detect actual environment - either "wsl" or "linux"
	if env != "wsl" && env != "linux" {
		t.Errorf("Expected environment to be 'wsl' or 'linux', got: %s", env)
	}
	
	// Should not return hardcoded value
	if env == "" {
		t.Errorf("Expected environment detection to return non-empty value")
	}
}

func TestShowEnvironmentSelection_ShouldUseDetectedEnvironmentAsDefault(t *testing.T) {
	cfg := config.Config{}
	tui := New(cfg)
	
	// This test expects that ShowEnvironmentSelection uses actual environment detection
	// as the default value, not hardcoded "linux"
	detected, err := tui.DetectActualEnvironment()
	if err != nil {
		t.Skip("Skipping test - environment detection failed")
	}
	
	selected, err := tui.ShowEnvironmentSelection()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// In a non-interactive test environment, it should return the detected environment
	// rather than always returning "linux"
	if detected == "wsl" && selected == "linux" {
		t.Errorf("Expected ShowEnvironmentSelection to use detected environment 'wsl' as default, got 'linux'")
	}
}

func TestShowToolSelection_ShouldReturnToolsFromAllCategories(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
				"zsh": config.ToolConfig{
					DisplayName: "Zsh Shell",
					BinaryName:  "zsh",
				},
			},
			"editors": {
				"vim": config.ToolConfig{
					DisplayName: "Vim Editor",
					BinaryName:  "vim",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	selectedTools, err := tui.ShowToolSelection()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Should return a non-empty slice when tools are available
	if len(selectedTools) == 0 {
		t.Errorf("Expected ShowToolSelection to return selected tools, got empty slice")
	}
	
	// Should be able to return tools from multiple categories
	// In a real scenario, this would be user-selected, but for testing
	// we expect it to potentially return tools from both categories
}

func TestShowToolSelection_ShouldUseCategorizedSelection(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	// ShowToolSelection should internally use ShowToolSelectionByCategory
	// and flatten the results into a simple []string
	categorySelections, err := tui.ShowToolSelectionByCategory()
	if err != nil {
		t.Skip("Skipping test - categorized selection failed")
	}
	
	toolSelections, err := tui.ShowToolSelection()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// The tool selection should be related to the categorized selection
	// At minimum, if we have categories, we should get some tool selection capability
	if len(categorySelections.CategoryAndTools) > 0 && len(toolSelections) == 0 {
		t.Errorf("Expected ShowToolSelection to return tools when categories are available, got empty")
	}
}

func TestCreateInteractiveToolForm_ShouldReturnStructuredForm(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
				"zsh": config.ToolConfig{
					DisplayName: "Zsh Shell",
					BinaryName:  "zsh",
				},
			},
			"editors": {
				"vim": config.ToolConfig{
					DisplayName: "Vim Editor",
					BinaryName:  "vim",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	// Should create a structured form for interactive selection
	formGroups, err := tui.CreateInteractiveToolForm()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if formGroups == nil {
		t.Errorf("Expected form groups to not be nil")
	}
	
	if len(formGroups) == 0 {
		t.Errorf("Expected form groups to contain category forms, got empty")
	}
	
	// Should have forms for each category
	if len(formGroups) != 2 {
		t.Errorf("Expected 2 form groups (shells, editors), got %d", len(formGroups))
	}
}

func TestShowInteractiveToolSelection_ShouldUseInteractiveForm(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	// This method should provide interactive selection using huh library
	selections, err := tui.ShowInteractiveToolSelection()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if selections.CategoryAndTools == nil {
		t.Errorf("Expected selections.CategoryAndTools to not be nil")
	}
	
	// Should return structured selections (even if empty in test environment)
	// The structure should match our existing Selections format
}

func TestRunInteractiveFormWithDefaults_ShouldExecuteFormLogic(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
				"zsh": config.ToolConfig{
					DisplayName: "Zsh Shell",
					BinaryName:  "zsh",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	// Should execute form logic and capture selections
	// In test environment, should use default values or skip interactive parts
	selections, err := tui.RunInteractiveFormWithDefaults()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if selections.CategoryAndTools == nil {
		t.Errorf("Expected selections.CategoryAndTools to not be nil")
	}
	
	// Should have processed the categories
	if len(selections.CategoryAndTools) == 0 {
		t.Errorf("Expected selections to contain categories, got empty")
	}
	
	// Should contain the shells category
	found := false
	for _, categoryTools := range selections.CategoryAndTools {
		if categoryTools.Category == "shells" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find 'shells' category in selections")
	}
}

func TestExecuteInteractiveForm_ShouldHandleFormExecution(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
				},
			},
		},
	}
	
	tui := New(cfg)
	
	// Create the form groups
	formGroups, err := tui.CreateInteractiveToolForm()
	if err != nil {
		t.Skip("Skipping test - form creation failed")
	}
	
	// Should execute the interactive form (or simulate execution in tests)
	selections, err := tui.ExecuteInteractiveForm(formGroups)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if selections.CategoryAndTools == nil {
		t.Errorf("Expected selections.CategoryAndTools to not be nil")
	}
	
	// Should return proper structure from form execution
}
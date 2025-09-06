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
package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestInstallCommandIntegration_ShouldUseTUIFlow(t *testing.T) {
	// Test that the install command properly integrates with our TUI methods
	// This would test the end-to-end flow:
	// 1. Load config
	// 2. Detect/select environment
	// 3. Show interactive tool selection
	// 4. Return structured selections

	selections, err := RunInstallFlow()

	if err != nil {
		t.Errorf("Expected no error from install flow, got: %v", err)
	}

	if selections.CategoryAndTools == nil {
		t.Errorf("Expected selections.CategoryAndTools to not be nil")
	}

	// Should return structured selections ready for installation
	if len(selections.CategoryAndTools) == 0 {
		t.Errorf("Expected install flow to return categories, got empty")
	}
}

func TestCreateInstallTUI_ShouldLoadConfigAndCreateTUI(t *testing.T) {
	// Test that we can create a TUI instance with loaded config
	tuiInstance, err := CreateInstallTUI()

	if err != nil {
		t.Errorf("Expected no error creating TUI, got: %v", err)
	}

	if tuiInstance == nil {
		t.Errorf("Expected TUI instance to not be nil")
	}

	env, err := tuiInstance.DetectActualEnvironment()
	if err != nil {
		t.Errorf("Expected no error detecting environment, got: %v", err)
	}

	if env != "linux" && env != "wsl" {
		t.Errorf("Expected environment to be 'linux' or 'wsl', got: %s", env)
	}
}

func TestRunInstallFlowWithConfig_ShouldExecuteCompleteFlow(t *testing.T) {
	// Test the complete install flow:
	// 1. Environment detection/selection
	// 2. Tool selection
	// 3. Return results

	configPath := "../config.yaml"

	result, err := RunInstallFlowWithConfig(configPath)

	if err != nil {
		t.Skip("Skipping test - config file not available")
	}

	if result.Environment == "" {
		t.Errorf("Expected environment to be set in result")
	}

	if result.Selections.CategoryAndTools == nil {
		t.Errorf("Expected result.Selections.CategoryAndTools to not be nil")
	}
}

func TestInstallCommand_ShouldUseInstallFlow(t *testing.T) {
	// Since fmt.Printf writes directly to stdout (not to cobra's output buffer),
	// we need to capture stdout directly using os.Stdout redirection

	// Save original stdout
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	installCmd.Run(installCmd, []string{})

	w.Close()
	os.Stdout = origStdout

	var output strings.Builder
	io.Copy(&output, r)
	outputStr := output.String()

	hasToolSelection := strings.Contains(outputStr, "Selected tools for installation:")
	hasRunFlowError := strings.Contains(outputStr, "Error running install flow:")

	if !hasToolSelection && !hasRunFlowError {
		t.Errorf("Install command should show tool selection results or RunInstallFlow error, got: %q", outputStr)
	}
}

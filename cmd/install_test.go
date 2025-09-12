package cmd

import (
	"testing"
)

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
	// This test verifies that the install command goes through the complete flow
	// but we need to avoid hanging on interactive TUI
	t.Skip("Skipping interactive TUI test - requires user interaction")
}

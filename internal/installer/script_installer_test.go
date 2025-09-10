package installer

import (
	"errors"
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
)

func TestScriptInstaller_ShouldExecuteInstallScript(t *testing.T) {
	// Test that ScriptInstaller executes the correct install script for a tool
	mockExecutor := &MockCommandExecutor{}
	installer := &ScriptInstaller{
		CommandExecutor: mockExecutor,
	}

	tool := config.ToolConfig{
		DisplayName:   "FZF Fuzzy Finder",
		BinaryName:    "fzf",
		InstallMethod: "script",
		PackageName:   "",
		InstallScript: "install_scripts/fzf.sh",
		Dependencies:  []string{"git"},
	}

	err := installer.Install(tool)

	if err != nil {
		t.Errorf("Expected ScriptInstaller to successfully install fzf, got error: %v", err)
	}

	// Should execute the install script
	expectedCmd := "bash install_scripts/fzf.sh"
	if len(mockExecutor.ExecutedCommands) != 1 {
		t.Errorf("Expected 1 command to be executed, got %d", len(mockExecutor.ExecutedCommands))
	}

	if len(mockExecutor.ExecutedCommands) > 0 && mockExecutor.ExecutedCommands[0] != expectedCmd {
		t.Errorf("Expected command '%s', got '%s'", expectedCmd, mockExecutor.ExecutedCommands[0])
	}
}

func TestScriptInstaller_ShouldHandleScriptFailure(t *testing.T) {
	// Test that ScriptInstaller properly handles script execution failures
	mockExecutor := &MockCommandExecutor{
		ShouldFail:   true,
		FailureError: errors.New("script execution failed: permission denied"),
	}
	installer := &ScriptInstaller{
		CommandExecutor: mockExecutor,
	}

	tool := config.ToolConfig{
		DisplayName:   "Broken Script Tool",
		BinaryName:    "broken-tool",
		InstallMethod: "script",
		InstallScript: "install_scripts/broken.sh",
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	// Should return an error when script execution fails
	if err == nil {
		t.Errorf("Expected ScriptInstaller to return error when script fails, got nil")
	}

	// Verify error contains context about the script failure
	if err != nil {
		expectedMsg := "failed to execute install script"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain '%s', got: %s", expectedMsg, err.Error())
		}
	}
}

func TestScriptInstaller_ShouldValidateScriptPath(t *testing.T) {
	// Test that ScriptInstaller validates that InstallScript is provided
	mockExecutor := &MockCommandExecutor{}
	installer := &ScriptInstaller{
		CommandExecutor: mockExecutor,
	}

	tool := config.ToolConfig{
		DisplayName:   "Tool Without Script",
		BinaryName:    "no-script-tool",
		InstallMethod: "script",
		InstallScript: "", // Empty script path should cause error
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	// Should return an error when no script is provided
	if err == nil {
		t.Errorf("Expected ScriptInstaller to return error when InstallScript is empty, got nil")
	}

	// Should not execute any commands when script path is invalid
	if len(mockExecutor.ExecutedCommands) > 0 {
		t.Errorf("Expected no commands to be executed when script path is empty, got %d commands", len(mockExecutor.ExecutedCommands))
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

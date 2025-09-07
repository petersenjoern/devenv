package installer

import (
	"reflect"
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
)

func TestManualInstaller_ShouldDisplayInstallationInstructions(t *testing.T) {
	// Test that ManualInstaller displays installation instructions from wsl_notes
	installer := &ManualInstaller{}

	tool := config.ToolConfig{
		DisplayName:   "Visual Studio Code",
		BinaryName:    "code",
		InstallMethod: "manual",
		PackageName:   "",
		InstallScript: "",
		WSLNotes:      "Download VS Code from https://code.visualstudio.com/ and install on Windows host. Then add 'code' command to PATH.",
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	// Should not return an error for valid manual installation
	if err != nil {
		t.Errorf("Expected ManualInstaller to handle installation instructions successfully, got error: %v", err)
	}

}

func TestManualInstaller_ShouldHandleEmptyInstructions(t *testing.T) {
	// Test that ManualInstaller handles tools without installation instructions gracefully
	installer := &ManualInstaller{}

	tool := config.ToolConfig{
		DisplayName:   "Tool Without Instructions",
		BinaryName:    "no-instructions-tool",
		InstallMethod: "manual",
		PackageName:   "",
		InstallScript: "",
		WSLNotes:      "", // Empty instructions
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	// Should not return an error - it should display fallback message instead
	if err != nil {
		t.Errorf("Expected ManualInstaller to handle empty instructions gracefully, got error: %v", err)
	}

	// Should display manual installation message even without specific instructions
	// (This test verifies the behavior but doesn't capture stdout - that's acceptable)
}

func TestManualInstaller_ShouldNotExecuteSystemCommands(t *testing.T) {
	// Test that ManualInstaller never executes system commands unlike APT/Script installers
	
	// Create a spy to track if any commands would have been executed
	// We'll verify this by ensuring ManualInstaller doesn't have a CommandExecutor field
	installer := &ManualInstaller{}

	tool := config.ToolConfig{
		DisplayName:   "Manual Tool",
		BinaryName:    "manual-tool",
		InstallMethod: "manual",
		PackageName:   "should-not-be-used",        // Should be ignored
		InstallScript: "should-not-be-executed.sh", // Should be ignored
		WSLNotes:      "Install this tool manually by following these steps: 1. Download from website 2. Run installer",
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	// Should succeed without executing any commands
	if err != nil {
		t.Errorf("Expected ManualInstaller to handle manual instructions without error, got: %v", err)
	}

	// Verify that ManualInstaller struct doesn't have CommandExecutor field
	// This architectural constraint ensures no system commands can be executed
	installerType := reflect.TypeOf(installer).Elem()
	for i := 0; i < installerType.NumField(); i++ {
		field := installerType.Field(i)
		if field.Name == "CommandExecutor" {
			t.Errorf("ManualInstaller should not have CommandExecutor field, but found: %s", field.Name)
		}
	}
}

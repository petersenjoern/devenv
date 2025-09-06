package installer

import (
	"errors"
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
)

type MockCommandExecutor struct {
	ExecutedCommands []string
	ShouldFail       bool
	FailureError     error
}

func (m *MockCommandExecutor) Execute(command string) error {
	m.ExecutedCommands = append(m.ExecutedCommands, command)
	if m.ShouldFail {
		return m.FailureError
	}
	return nil
}

func TestAPTInstaller_ShouldExecuteCorrectCommands(t *testing.T) {
	// Test that APTInstaller executes the correct apt commands without affecting the OS
	mockExecutor := &MockCommandExecutor{}
	installer := &APTInstaller{
		CommandExecutor: mockExecutor,
	}

	tool := config.ToolConfig{
		DisplayName:   "Git Version Control",
		BinaryName:    "git",
		InstallMethod: "apt",
		PackageName:   "git",
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	if err != nil {
		t.Errorf("Expected APTInstaller to successfully install git, got error: %v", err)
	}

	expectedCommands := []string{
		"sudo apt update",
		"sudo apt install -y git",
	}

	if len(mockExecutor.ExecutedCommands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(mockExecutor.ExecutedCommands))
	}

	for i, expected := range expectedCommands {
		if i >= len(mockExecutor.ExecutedCommands) {
			t.Errorf("Missing command at index %d: expected %s", i, expected)
			continue
		}
		if mockExecutor.ExecutedCommands[i] != expected {
			t.Errorf("Command at index %d: expected %s, got %s", i, expected, mockExecutor.ExecutedCommands[i])
		}
	}
}

func TestAPTInstaller_ShouldHandleCommandFailure(t *testing.T) {
	// Test that APTInstaller properly handles command execution failures
	mockExecutor := &MockCommandExecutor{
		ShouldFail:   true,
		FailureError: errors.New("apt command failed: package not found"),
	}
	installer := &APTInstaller{
		CommandExecutor: mockExecutor,
	}

	tool := config.ToolConfig{
		DisplayName:   "Non-existent Package",
		BinaryName:    "nonexistent-tool",
		InstallMethod: "apt",
		PackageName:   "nonexistent-package",
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	if err == nil {
		t.Errorf("Expected APTInstaller to return error when command fails, got nil")
	}

	if err != nil && !errors.Is(err, mockExecutor.FailureError) {
		// Check if the error is wrapped properly
		if err.Error() == "" || len(err.Error()) == 0 {
			t.Errorf("Expected error message to contain context, got empty error")
		}
		// The error should be from apt update (first command) in this test
		if len(mockExecutor.ExecutedCommands) > 0 && mockExecutor.ExecutedCommands[0] != "sudo apt update" {
			t.Errorf("Expected first command to be apt update, got: %s", mockExecutor.ExecutedCommands[0])
		}
	}
}

func TestAPTInstaller_ShouldRunUpdateBeforeInstall(t *testing.T) {
	// Test command order - apt update should run before apt install
	mockExecutor := &MockCommandExecutor{}
	installer := &APTInstaller{
		CommandExecutor: mockExecutor,
	}

	tool := config.ToolConfig{
		DisplayName:   "Vim Editor",
		BinaryName:    "vim",
		InstallMethod: "apt",
		PackageName:   "vim",
		Dependencies:  []string{},
	}

	err := installer.Install(tool)

	if err != nil {
		t.Errorf("Expected APTInstaller to successfully install vim, got error: %v", err)
	}

	// Verify that exactly 2 commands were executed
	expectedCommands := 2
	if len(mockExecutor.ExecutedCommands) != expectedCommands {
		t.Errorf("Expected %d commands to be executed, got %d", expectedCommands, len(mockExecutor.ExecutedCommands))
	}

	// Verify command order
	if len(mockExecutor.ExecutedCommands) >= 1 {
		if mockExecutor.ExecutedCommands[0] != "sudo apt update" {
			t.Errorf("Expected first command to be 'sudo apt update', got: %s", mockExecutor.ExecutedCommands[0])
		}
	}

	if len(mockExecutor.ExecutedCommands) >= 2 {
		expectedInstallCmd := "sudo apt install -y vim"
		if mockExecutor.ExecutedCommands[1] != expectedInstallCmd {
			t.Errorf("Expected second command to be '%s', got: %s", expectedInstallCmd, mockExecutor.ExecutedCommands[1])
		}
	}
}

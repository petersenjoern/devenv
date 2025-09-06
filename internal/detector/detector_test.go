package detector

import (
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
)

func TestIsBinaryInstalled_ShouldReturnTrueForInstalledBinary(t *testing.T) {
	detector := New()

	// 'ls' should always be available on Ubuntu/Linux systems
	installed := detector.IsBinaryInstalled("ls")

	if !installed {
		t.Errorf("Expected 'ls' to be installed, got false")
	}
}

func TestIsBinaryInstalled_ShouldReturnFalseForNonExistentBinary(t *testing.T) {
	detector := New()

	// 'nonexistent-binary-12345' should not exist
	installed := detector.IsBinaryInstalled("nonexistent-binary-12345")

	if installed {
		t.Errorf("Expected 'nonexistent-binary-12345' to not be installed, got true")
	}
}

func TestDetectTool_ShouldDetectInstalledBinary(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "ls",
	}

	status := detector.DetectTool(tool)

	if !status.BinaryInstalled {
		t.Errorf("Expected BinaryInstalled to be true for 'ls', got false")
	}
}

func TestDetectTool_ShouldDetectNonExistentBinary(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "nonexistent-binary-12345",
	}

	status := detector.DetectTool(tool)

	if status.BinaryInstalled {
		t.Errorf("Expected BinaryInstalled to be false for 'nonexistent-binary-12345', got true")
	}
}

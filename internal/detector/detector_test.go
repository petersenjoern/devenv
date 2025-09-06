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

func TestDetectTool_ShouldReturnPathForInstalledBinary(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "ls",
	}

	status := detector.DetectTool(tool)

	if status.Path == "" {
		t.Errorf("Expected Path to be set for installed binary 'ls', got empty string")
	}

	if status.Path != "/usr/bin/ls" && status.Path != "/bin/ls" {
		t.Logf("Path for 'ls' is: %s (this may vary by system)", status.Path)
	}
}

func TestDetectTool_ShouldReturnEmptyPathForNonExistentBinary(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "nonexistent-binary-12345",
	}

	status := detector.DetectTool(tool)

	if status.Path != "" {
		t.Errorf("Expected Path to be empty for non-existent binary, got '%s'", status.Path)
	}
}

func TestDetectTool_ShouldReturnVersionForInstalledBinary(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "bash",
	}

	status := detector.DetectTool(tool)

	if status.Version == "" {
		t.Errorf("Expected Version to be set for installed binary 'bash', got empty string")
	}
}

func TestDetectTool_ShouldReturnEmptyVersionForNonExistentBinary(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "nonexistent-binary-12345",
	}

	status := detector.DetectTool(tool)

	if status.Version != "" {
		t.Errorf("Expected Version to be empty for non-existent binary, got '%s'", status.Version)
	}
}

func TestDetectTool_ShouldDetectConfigWhenFileExists(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "bash",
		ConfigPath: "/etc/passwd", // This file should exist on all Linux systems
	}

	status := detector.DetectTool(tool)

	if !status.ConfigApplied {
		t.Errorf("Expected ConfigApplied to be true when config file exists, got false")
	}
}

func TestDetectTool_ShouldReturnFalseConfigWhenFileDoesNotExist(t *testing.T) {
	detector := New()

	tool := config.ToolConfig{
		BinaryName: "bash",
		ConfigPath: "/nonexistent/config/path/file.conf",
	}

	status := detector.DetectTool(tool)

	if status.ConfigApplied {
		t.Errorf("Expected ConfigApplied to be false when config file doesn't exist, got true")
	}
}

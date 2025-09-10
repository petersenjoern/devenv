package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/installer"
	"github.com/petersenjoern/devenv/internal/tui"
)

func TestInstallCommand_ShouldFindConfigPath(t *testing.T) {
	// Test that install command can find config file using standard paths

	configPath, err := findConfigPath()

	if err != nil {
		t.Skipf("No config file found for testing, skipping: %v", err)
		return
	}

	if configPath == "" {
		t.Errorf("Expected non-empty config path")
	}

	expectedPaths := []string{"./config.yaml", "../config.yaml"}
	validPath := false
	for _, expectedPath := range expectedPaths {
		if configPath == expectedPath {
			validPath = true
			break
		}
	}

	if !validPath {
		t.Errorf("Expected config path to be one of %v, got %s", expectedPaths, configPath)
	}
}

func TestInstallCommand_ShouldDisplayInstallationResults(t *testing.T) {
	// Test that install command displays installation results properly

	// Mock installation results with mixed success/failure
	mockResults := map[string]installer.InstallationResult{
		"git": {
			Tool: config.ToolConfig{
				DisplayName: "Git Version Control",
				BinaryName:  "git",
			},
			Success: true,
			Error:   nil,
		},
		"nonexistent-tool": {
			Tool: config.ToolConfig{
				DisplayName: "Nonexistent Tool",
				BinaryName:  "nonexistent-tool",
			},
			Success: false,
			Error:   fmt.Errorf("installation failed: package not found"),
		},
	}

	// Capture output
	var output strings.Builder
	originalOutput := captureOutput(&output)

	displayInstallationResults(mockResults)

	originalOutput.restore()
	outputStr := output.String()

	// Should display installation results header
	if !strings.Contains(outputStr, "Installation Results") {
		t.Errorf("Expected output to contain 'Installation Results', got: %s", outputStr)
	}

	// Should show successful installation
	if !strings.Contains(outputStr, "✓") && !strings.Contains(outputStr, "Git Version Control") {
		t.Errorf("Expected output to show successful git installation, got: %s", outputStr)
	}

	// Should show failed installation
	if !strings.Contains(outputStr, "✗") && !strings.Contains(outputStr, "Nonexistent Tool") {
		t.Errorf("Expected output to show failed installation, got: %s", outputStr)
	}

	// Should display summary
	if !strings.Contains(outputStr, "Summary") {
		t.Errorf("Expected output to contain 'Summary', got: %s", outputStr)
	}

	// Should show counts
	if !strings.Contains(outputStr, "Total attempted: 2") {
		t.Errorf("Expected output to show total attempted count, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "Successful: 1") {
		t.Errorf("Expected output to show successful count, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "Failed: 1") {
		t.Errorf("Expected output to show failed count, got: %s", outputStr)
	}
}

func TestInstallCommand_ShouldHandleNoToolsSelected(t *testing.T) {
	// Test that install command handles case where user selects no tools

	emptySelections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{},
	}

	configPath, err := findConfigPath()
	if err != nil {
		t.Skipf("No config file found for testing, skipping: %v", err)
		return
	}

	// Execute installations with empty selections
	results, err := ExecuteInstallations(emptySelections, configPath)

	// Should not return error even with no selections
	if err != nil {
		t.Errorf("Expected no error with empty selections, got: %v", err)
	}

	// Should return empty results map
	if len(results) != 0 {
		t.Errorf("Expected empty results with no selections, got %d results", len(results))
	}

	// Test output with empty results
	var output strings.Builder
	originalOutput := captureOutput(&output)

	displayInstallationResults(results)

	originalOutput.restore()
	outputStr := output.String()

	// Should still show headers but with zero counts
	if !strings.Contains(outputStr, "Installation Results") {
		t.Errorf("Expected output to contain 'Installation Results' even with no tools, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "Total attempted: 0") {
		t.Errorf("Expected output to show zero attempted installations, got: %s", outputStr)
	}
}

func TestInstallCommand_ShouldCountSelectedTools(t *testing.T) {
	// Test that install command correctly counts selected tools across categories

	mockSelections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{
			{
				Category: "version_control",
				Tools:    []string{"git", "mercurial"},
			},
			{
				Category: "editors",
				Tools:    []string{"vim", "emacs", "nano"},
			},
		},
	}

	// Use mock tool configs instead of loading from real config
	mockToolConfigs := map[string]config.ToolConfig{
		"git": {
			DisplayName:   "Git Version Control",
			BinaryName:    "git",
			InstallMethod: "apt",
			PackageName:   "git",
		},
		"mercurial": {
			DisplayName:   "Mercurial Version Control",
			BinaryName:    "hg",
			InstallMethod: "apt",
			PackageName:   "mercurial",
		},
		"vim": {
			DisplayName:   "Vim Editor",
			BinaryName:    "vim",
			InstallMethod: "apt",
			PackageName:   "vim",
		},
		"emacs": {
			DisplayName:   "Emacs Editor",
			BinaryName:    "emacs",
			InstallMethod: "apt",
			PackageName:   "emacs",
		},
		"nano": {
			DisplayName:   "Nano Editor",
			BinaryName:    "nano",
			InstallMethod: "apt",
			PackageName:   "nano",
		},
	}

	orchestrator := &installer.InstallationOrchestrator{
		APTInstaller:    &installer.APTInstaller{CommandExecutor: &MockCommandExecutor{}},
		ScriptInstaller: &installer.ScriptInstaller{CommandExecutor: &MockCommandExecutor{}},
		ManualInstaller: &installer.ManualInstaller{},
	}

	// Execute installations with mocked orchestrator
	results := orchestrator.ExecuteInstallations(mockSelections, mockToolConfigs)

	// Test that the actual implementation counts tools correctly
	expectedTotal := 5 // 2 version control + 3 editors
	actualTotal := len(results)

	if actualTotal != expectedTotal {
		t.Errorf("Expected %d tools in results, got %d", expectedTotal, actualTotal)
	}

	// Test the display output shows correct count
	var output strings.Builder
	originalOutput := captureOutput(&output)

	displayInstallationResults(results)

	originalOutput.restore()
	outputStr := output.String()

	// Should show the correct total in summary
	expectedSummary := fmt.Sprintf("Total attempted: %d", expectedTotal)
	if !strings.Contains(outputStr, expectedSummary) {
		t.Errorf("Expected output to contain '%s', got: %s", expectedSummary, outputStr)
	}
}

func TestInstallCommand_ShouldProvideHelpfulGuidanceOnFailures(t *testing.T) {
	// Test that install command provides helpful next steps when installations fail

	// Mock results with some failures
	mockResultsWithFailures := map[string]installer.InstallationResult{
		"git": {
			Tool: config.ToolConfig{
				DisplayName: "Git Version Control",
				BinaryName:  "git",
			},
			Success: true,
			Error:   nil,
		},
		"failed-tool": {
			Tool: config.ToolConfig{
				DisplayName: "Failed Tool",
				BinaryName:  "failed-tool",
			},
			Success: false,
			Error:   fmt.Errorf("installation failed"),
		},
	}

	// Capture output
	var output strings.Builder
	originalOutput := captureOutput(&output)

	displayInstallationResults(mockResultsWithFailures)

	originalOutput.restore()
	outputStr := output.String()

	// Should provide helpful guidance for failures
	if !strings.Contains(outputStr, "Some installations failed") {
		t.Errorf("Expected guidance about failed installations, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "devenv status") {
		t.Errorf("Expected guidance to check status, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "devenv install") {
		t.Errorf("Expected guidance to retry installation, got: %s", outputStr)
	}
}

func TestInstallCommand_ShouldProvideSuccessGuidanceOnAllSuccess(t *testing.T) {
	// Test that install command provides success guidance when all installations succeed

	// Mock results with all successes
	mockResultsAllSuccess := map[string]installer.InstallationResult{
		"git": {
			Tool: config.ToolConfig{
				DisplayName: "Git Version Control",
				BinaryName:  "git",
			},
			Success: true,
			Error:   nil,
		},
		"vim": {
			Tool: config.ToolConfig{
				DisplayName: "Vim Editor",
				BinaryName:  "vim",
			},
			Success: true,
			Error:   nil,
		},
	}

	// Capture output
	var output strings.Builder
	originalOutput := captureOutput(&output)

	displayInstallationResults(mockResultsAllSuccess)

	originalOutput.restore()
	outputStr := output.String()

	// Should provide success guidance
	if !strings.Contains(outputStr, "All installations completed successfully") {
		t.Errorf("Expected success message, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "devenv status") {
		t.Errorf("Expected guidance to verify environment, got: %s", outputStr)
	}

	// Should not show failure guidance
	if strings.Contains(outputStr, "Some installations failed") {
		t.Errorf("Should not show failure guidance on all success, got: %s", outputStr)
	}
}

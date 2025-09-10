package installer

import (
	"fmt"
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/tui"
)

func TestOrchestrator_ShouldRouteToolsToCorrectInstallers(t *testing.T) {
	// Test that orchestrator uses the correct installer based on install_method
	mockAPTExecutor := &MockCommandExecutor{}
	mockScriptExecutor := &MockCommandExecutor{}

	orchestrator := &InstallationOrchestrator{
		APTInstaller:    &APTInstaller{CommandExecutor: mockAPTExecutor},
		ScriptInstaller: &ScriptInstaller{CommandExecutor: mockScriptExecutor},
		ManualInstaller: &ManualInstaller{},
	}

	selections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{
			{
				Category: "utilities",
				Tools:    []string{"git", "docker"}, // git=apt, docker=script
			},
		},
	}

	tools := map[string]config.ToolConfig{
		"git": {
			DisplayName:   "Git Version Control",
			BinaryName:    "git",
			InstallMethod: "apt",
			PackageName:   "git",
			Dependencies:  []string{},
		},
		"docker": {
			DisplayName:   "Docker Engine",
			BinaryName:    "docker",
			InstallMethod: "script",
			InstallScript: "install_scripts/docker.sh",
			Dependencies:  []string{"wget"},
		},
	}

	results := orchestrator.ExecuteInstallations(selections, tools)

	// Should have results for both tools
	if len(results) != 2 {
		t.Errorf("Expected 2 installation results, got %d", len(results))
	}

	// Should have used APT installer for git
	expectedAPTCommands := []string{"sudo apt update", "sudo apt install -y git"}
	if len(mockAPTExecutor.ExecutedCommands) != len(expectedAPTCommands) {
		t.Errorf("Expected %d APT commands, got %d", len(expectedAPTCommands), len(mockAPTExecutor.ExecutedCommands))
	}

	// Should have used Script installer for docker
	expectedScriptCommands := []string{"bash install_scripts/docker.sh"}
	if len(mockScriptExecutor.ExecutedCommands) != len(expectedScriptCommands) {
		t.Errorf("Expected %d script commands, got %d", len(expectedScriptCommands), len(mockScriptExecutor.ExecutedCommands))
	}
}

func TestOrchestrator_ShouldHandleDependencyOrdering(t *testing.T) {
	// Test that dependencies are installed before dependent tools
	mockExecutor := &MockCommandExecutor{}

	orchestrator := &InstallationOrchestrator{
		APTInstaller:    &APTInstaller{CommandExecutor: mockExecutor},
		ScriptInstaller: &ScriptInstaller{CommandExecutor: mockExecutor},
		ManualInstaller: &ManualInstaller{},
	}

	selections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{
			{
				Category: "utilities",
				Tools:    []string{"lazydocker"}, // lazydocker depends on docker and curl
			},
		},
	}

	tools := map[string]config.ToolConfig{
		"curl": {
			DisplayName:   "Curl HTTP Client",
			BinaryName:    "curl",
			InstallMethod: "apt",
			PackageName:   "curl",
			Dependencies:  []string{},
		},
		"docker": {
			DisplayName:   "Docker Engine",
			BinaryName:    "docker",
			InstallMethod: "script",
			InstallScript: "install_scripts/docker.sh",
			Dependencies:  []string{},
		},
		"lazydocker": {
			DisplayName:   "Lazydocker Terminal UI",
			BinaryName:    "lazydocker",
			InstallMethod: "script",
			InstallScript: "install_scripts/lazydocker.sh",
			Dependencies:  []string{"docker", "curl"}, // depends on docker and curl
		},
	}

	results := orchestrator.ExecuteInstallations(selections, tools)

	// Should install curl first (apt), then docker (script), then lazydocker (script)
	if len(results) != 3 {
		t.Errorf("Expected 3 installation results (curl, docker, lazydocker), got %d", len(results))
	}

	// Should install curl first, then docker, then lazydocker
	//
	curlUpdateIndex := -1
	dockerInstallIndex := -1
	lazydockerInstallIndex := -1

	for i, cmd := range mockExecutor.ExecutedCommands {
		if cmd == "sudo apt update" {
			curlUpdateIndex = i
		}
		if cmd == "bash install_scripts/docker.sh" {
			dockerInstallIndex = i
		}
		if cmd == "bash install_scripts/lazydocker.sh" {
			lazydockerInstallIndex = i
		}
	}

	if curlUpdateIndex == -1 || dockerInstallIndex == -1 || lazydockerInstallIndex == -1 {
		t.Errorf("Missing expected commands in execution order")
	}

	if dockerInstallIndex >= lazydockerInstallIndex {
		t.Errorf("Docker should be installed before lazydocker (dependency ordering)")
	}
	if curlUpdateIndex >= lazydockerInstallIndex {
		t.Errorf("Curl should be installed before lazydocker (dependency ordering)")
	}
	if curlUpdateIndex >= dockerInstallIndex {
		t.Errorf("Curl should be installed before docker (dependency ordering)")
	}

	// Verify the exact command order
	expectedCommands := []string{
		"sudo apt update",                    // curl install
		"sudo apt install -y curl",           // curl install
		"bash install_scripts/docker.sh",     // docker install
		"bash install_scripts/lazydocker.sh", // lazydocker install
	}
	if len(mockExecutor.ExecutedCommands) != len(expectedCommands) {
		t.Errorf("Expected %d total commands, got %d", len(expectedCommands), len(mockExecutor.ExecutedCommands))
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

func TestOrchestrator_ShouldCollectInstallationResults(t *testing.T) {
	// Test that orchestrator collects and returns installation results
	mockExecutor := &MockCommandExecutor{}

	orchestrator := &InstallationOrchestrator{
		APTInstaller:    &APTInstaller{CommandExecutor: mockExecutor},
		ScriptInstaller: &ScriptInstaller{CommandExecutor: mockExecutor},
		ManualInstaller: &ManualInstaller{},
	}

	selections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{
			{
				Category: "utilities",
				Tools:    []string{"git"},
			},
		},
	}

	tools := map[string]config.ToolConfig{
		"git": {
			DisplayName:   "Git Version Control",
			BinaryName:    "git",
			InstallMethod: "apt",
			PackageName:   "git",
			Dependencies:  []string{},
		},
	}

	results := orchestrator.ExecuteInstallations(selections, tools)

	// Should have result for git
	gitResult, found := results["git"]
	if !found {
		t.Errorf("Expected result for git installation")
	}

	// Should indicate successful installation
	if gitResult.Success != true {
		t.Errorf("Expected git installation to succeed")
	}

	if gitResult.Tool.BinaryName != "git" {
		t.Errorf("Expected result tool to be git, got %s", gitResult.Tool.BinaryName)
	}
}

func TestOrchestrator_ShouldHandleInstallationFailures(t *testing.T) {
	// Test that orchestrator handles failures gracefully and continues with other tools
	mockExecutor := &MockCommandExecutor{
		ShouldFail:   true,
		FailureError: fmt.Errorf("apt command failed"),
	}

	orchestrator := &InstallationOrchestrator{
		APTInstaller:    &APTInstaller{CommandExecutor: mockExecutor},
		ScriptInstaller: &ScriptInstaller{CommandExecutor: mockExecutor},
		ManualInstaller: &ManualInstaller{},
	}

	selections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{
			{
				Category: "utilities",
				Tools:    []string{"git"},
			},
		},
	}

	tools := map[string]config.ToolConfig{
		"git": {
			DisplayName:   "Git Version Control",
			BinaryName:    "git",
			InstallMethod: "apt",
			PackageName:   "git",
			Dependencies:  []string{},
		},
	}

	results := orchestrator.ExecuteInstallations(selections, tools)

	// Should have result for git even though it failed
	gitResult, found := results["git"]
	if !found {
		t.Errorf("Expected result for git installation even on failure")
	}

	// Should indicate failed installation
	if gitResult.Success != false {
		t.Errorf("Expected git installation to fail")
	}

	// Should include error information
	if gitResult.Error == nil {
		t.Errorf("Expected error information in failed result")
	}
}

func TestOrchestrator_ShouldHandleManualInstallations(t *testing.T) {
	// Test that manual installations are handled (display messages, always succeed)
	orchestrator := &InstallationOrchestrator{
		APTInstaller:    &APTInstaller{CommandExecutor: &MockCommandExecutor{}},
		ScriptInstaller: &ScriptInstaller{CommandExecutor: &MockCommandExecutor{}},
		ManualInstaller: &ManualInstaller{},
	}

	selections := tui.Selections{
		CategoryAndTools: []tui.CategoryAndTools{
			{
				Category: "terminals",
				Tools:    []string{"alacritty"},
			},
		},
	}

	tools := map[string]config.ToolConfig{
		"alacritty": {
			DisplayName:   "Alacritty Terminal",
			BinaryName:    "alacritty",
			InstallMethod: "manual",
			WSLNotes:      "Install Alacritty on Windows host system.",
			Dependencies:  []string{},
		},
	}

	results := orchestrator.ExecuteInstallations(selections, tools)

	// Should have result for alacritty
	alacrittyResult, found := results["alacritty"]
	if !found {
		t.Errorf("Expected result for alacritty manual installation")
	}

	// Manual installations should always succeed (they just display instructions)
	if alacrittyResult.Success != true {
		t.Errorf("Expected manual installation to succeed")
	}
}

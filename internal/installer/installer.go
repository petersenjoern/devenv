package installer

import (
	"fmt"
	"os/exec"
	"sort"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/tui"
)

type CommandExecutor interface {
	Execute(command string) error
}

type RealCommandExecutor struct{}

func (r *RealCommandExecutor) Execute(command string) error {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Run()
}

type Installer interface {
	Install(tool config.ToolConfig) error
}

type APTInstaller struct {
	CommandExecutor CommandExecutor
}

type ScriptInstaller struct {
	CommandExecutor CommandExecutor
}

type ManualInstaller struct{}

func NewAPTInstaller() *APTInstaller {
	return &APTInstaller{
		CommandExecutor: &RealCommandExecutor{},
	}
}

func NewScriptInstaller() *ScriptInstaller {
	return &ScriptInstaller{
		CommandExecutor: &RealCommandExecutor{},
	}
}

const (
	aptUpdateCmd     = "sudo apt update"
	aptInstallCmd    = "sudo apt install -y %s"
	scriptInstallCmd = "bash %s"

	manualInstallMsg      = "Manual installation required for %s (%s)"
	manualInstructionsMsg = "Installation instructions:\n%s"
	manualFallbackMsg     = "No specific installation instructions provided. Please install %s manually."
	manualVerifyMsg       = "Please complete the installation manually and run 'devenv status' to verify."
)

func (a *APTInstaller) Install(tool config.ToolConfig) error {
	if err := a.CommandExecutor.Execute(aptUpdateCmd); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	installCmd := fmt.Sprintf(aptInstallCmd, tool.PackageName)
	if err := a.CommandExecutor.Execute(installCmd); err != nil {
		return fmt.Errorf("failed to install package %s: %w", tool.PackageName, err)
	}

	return nil
}

func (s *ScriptInstaller) Install(tool config.ToolConfig) error {
	if tool.InstallScript == "" {
		return fmt.Errorf("install script path is required for script installation method")
	}

	scriptCmd := fmt.Sprintf(scriptInstallCmd, tool.InstallScript)
	if err := s.CommandExecutor.Execute(scriptCmd); err != nil {
		return fmt.Errorf("failed to execute install script %s: %w", tool.InstallScript, err)
	}

	return nil
}

func (m *ManualInstaller) Install(tool config.ToolConfig) error {
	fmt.Printf(manualInstallMsg+"\n", tool.DisplayName, tool.BinaryName)

	if tool.WSLNotes != "" {
		fmt.Printf(manualInstructionsMsg+"\n", tool.WSLNotes)
	} else {
		fmt.Printf(manualFallbackMsg+"\n", tool.DisplayName)
	}

	fmt.Println(manualVerifyMsg)
	return nil
}

type InstallationOrchestrator struct {
	APTInstaller    *APTInstaller
	ScriptInstaller *ScriptInstaller
	ManualInstaller *ManualInstaller
}

type InstallationResult struct {
	Tool    config.ToolConfig
	Success bool
	Error   error
}

func (o *InstallationOrchestrator) ExecuteInstallations(selections tui.Selections, tools map[string]config.ToolConfig) map[string]InstallationResult {
	results := make(map[string]InstallationResult)

	selectedTools := o.extractSelectedTools(selections)

	installOrder := o.resolveDependencies(selectedTools, tools)

	for _, toolName := range installOrder {
		tool := tools[toolName]
		result := o.installTool(tool)
		results[toolName] = result
	}

	return results
}

func (o *InstallationOrchestrator) extractSelectedTools(selections tui.Selections) []string {
	var selectedTools []string
	for _, categoryAndTools := range selections.CategoryAndTools {
		selectedTools = append(selectedTools, categoryAndTools.Tools...)
	}
	return selectedTools
}

func (o *InstallationOrchestrator) resolveDependencies(selectedTools []string, tools map[string]config.ToolConfig) []string {
	// Build dependency graph and resolve installation order using topological sort
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var installOrder []string

	var visit func(toolName string)
	visit = func(toolName string) {
		if visited[toolName] || visiting[toolName] {
			return
		}

		visiting[toolName] = true

		if tool, exists := tools[toolName]; exists {
			deps := make([]string, len(tool.Dependencies))
			copy(deps, tool.Dependencies)
			sort.Strings(deps)

			for _, dep := range deps {
				if _, depExists := tools[dep]; depExists {
					visit(dep)
				}
			}
		}

		visiting[toolName] = false
		visited[toolName] = true
		installOrder = append(installOrder, toolName)
	}

	for _, toolName := range selectedTools {
		visit(toolName)
	}

	return installOrder
}

func (o *InstallationOrchestrator) installTool(tool config.ToolConfig) InstallationResult {
	var err error

	switch tool.InstallMethod {
	case "apt":
		err = o.APTInstaller.Install(tool)
	case "script":
		err = o.ScriptInstaller.Install(tool)
	case "manual":
		err = o.ManualInstaller.Install(tool)
	default:
		err = fmt.Errorf("unknown install method: %s", tool.InstallMethod)
	}

	return InstallationResult{
		Tool:    tool,
		Success: err == nil,
		Error:   err,
	}
}

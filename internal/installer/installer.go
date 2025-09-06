package installer

import (
	"fmt"
	"os/exec"

	"github.com/petersenjoern/devenv/internal/config"
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
	aptUpdateCmd    = "sudo apt update"
	aptInstallCmd   = "sudo apt install -y %s"
	scriptInstallCmd = "bash %s"
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
	return nil
}

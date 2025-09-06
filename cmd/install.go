package cmd

import (
	"fmt"
	"os"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devenv",
	Short: "DevEnv - Automated developer environment setup",
	Long: `DevEnv is a comprehensive tool for automating the setup of
development environments. It provides an interactive interface to
install and configure various development tools, programming languages,
and essential utilities needed for software development.`,
	Version: "0.1.0",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Launch interactive installation of development tools",
	Long: `Launch interactive TUI for tool selection and installation
First prompts for environment selection (WSL/Linux), 
then displays categorized tool selection with dependency resolution.`,
	Run: func(cmd *cobra.Command, args []string) {
		selections, err := RunInstallFlow()
		if err != nil {
			fmt.Printf("Error running install flow: %v\n", err)
			return
		}
		fmt.Printf("Selected tools for installation: %v\n", selections)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

type InstallResult struct {
	Environment string
	Selections  tui.Selections
}

func RunInstallFlow() (tui.Selections, error) {
	tuiInstance, err := CreateInstallTUI()
	if err != nil {
		return tui.Selections{}, err
	}

	return tuiInstance.ShowInteractiveToolSelection()
}

func CreateInstallTUI() (*tui.TUI, error) {
	// CLI entry vs unittest entry
	configPaths := []string{"./config.yaml", "../config.yaml"}

	for _, configPath := range configPaths {
		if _, err := os.Stat(configPath); err == nil {
			return createTUIFromConfig(configPath)
		}
	}

	return nil, fmt.Errorf("config file not found, tried: %v", configPaths)
}

// potentially delete this function as it's not used anywhere else
func RunInstallFlowWithConfig(configPath string) (InstallResult, error) {
	var result InstallResult

	tuiInstance, err := createTUIFromConfig(configPath)
	if err != nil {
		return result, fmt.Errorf("failed to create TUI instance: %w", err)
	}

	env, err := tuiInstance.DetectActualEnvironment()
	if err != nil {
		return result, fmt.Errorf("failed to detect environment: %w", err)
	}
	result.Environment = env

	selections, err := tuiInstance.ShowInteractiveToolSelection()
	if err != nil {
		return result, fmt.Errorf("failed to run interactive tool selection: %w", err)
	}
	result.Selections = selections

	return result, nil
}

func createTUIFromConfig(configPath string) (*tui.TUI, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	return tui.New(cfg), nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}

package cmd

import (
	"fmt"

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
		fmt.Println("Install command - TUI implementation coming soon")

		SelectedCateogiresAndTools, err := tui.CreateInstallationForm()
		if err != nil {
			fmt.Printf("Error creating installation form: %v\n", err)
			return
		}
		fmt.Printf("Selected tools for installation: %v\n", SelectedCateogiresAndTools)
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
	return createTUIFromConfig("./config.yaml")
}

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

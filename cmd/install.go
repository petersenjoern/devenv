package cmd

import (
	"fmt"
	"os"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/detector"
	"github.com/petersenjoern/devenv/internal/installer"
	"github.com/petersenjoern/devenv/internal/tui"
	"github.com/spf13/cobra"
)

const (
	resultsHeader = "\n=== Installation Results ==="
	summaryHeader = "\n=== Summary ==="
	successIcon   = "✓"
	failureIcon   = "✗"
	successMsg    = "All installations completed successfully!"
	failureMsg    = "Some installations failed. You can:"
	statusCmdStr  = "devenv status"
	retryCmd      = "devenv install"
)

var defaultConfigsPaths = []string{"./config.yaml", "../config.yaml"}

var rootCmd = &cobra.Command{
	Use:   "devenv",
	Short: "DevEnv - Automated developer environment setup",
	Long: `DevEnv is a comprehensive tool for automating the setup of
development environments. It provides an interactive interface to
install and configure various development tools and utilities.`,
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

		installResult, err := RunInstallFlowWithConfig("./config.yaml")
		if err != nil {
			fmt.Printf("Error running install flow: %v\n", err)
			return
		}
		fmt.Printf("Detected environment: %s\n", installResult.Environment)

		configPath, err := findConfigPath()
		if err != nil {
			fmt.Printf("Error finding config: %v\n", err)
			return
		}

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		detector := detector.New()
		toolsWithStatus := GetAllToolsStatus(cfg, detector)

		results, err := ExecuteInstallations(installResult.Selections, cfg, toolsWithStatus)
		if err != nil {
			fmt.Printf("Error executing installations: %v\n", err)
			return
		}

		displayInstallationResults(results)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

type InstallResult struct {
	Environment string
	Selections  tui.Selections
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

func ExecuteInstallations(selections tui.Selections, cfg config.Config, toolsWithStatus map[string]map[detector.Status]config.ToolConfig) (map[string]installer.InstallationResult, error) {

	toolConfigs := configToToolConfig(cfg)
	orchestrator := CreateInstallationOrchestrator()
	alreadyInstalledTools := []string{}
	for toolName, statusMap := range toolsWithStatus {
		for status := range statusMap {
			if status.BinaryInstalled == true {
				alreadyInstalledTools = append(alreadyInstalledTools, toolName)
			}
		}
	}

	results := orchestrator.ExecuteInstallations(selections, toolConfigs, alreadyInstalledTools)

	return results, nil
}

func configToToolConfig(cfg config.Config) map[string]config.ToolConfig {
	toolConfigs := make(map[string]config.ToolConfig)
	for _, category := range cfg.Categories {
		for toolName, toolConfig := range category {
			toolConfigs[toolName] = toolConfig
		}
	}
	return toolConfigs
}

func CreateInstallationOrchestrator() *installer.InstallationOrchestrator {
	return &installer.InstallationOrchestrator{
		APTInstaller:    installer.NewAPTInstaller(),    // Real APT command execution
		ScriptInstaller: installer.NewScriptInstaller(), // Real script execution
		ManualInstaller: &installer.ManualInstaller{},   // User instruction display
	}
}

func findConfigPath() (string, error) {

	for _, configPath := range defaultConfigsPaths {
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	return "", fmt.Errorf("config file not found, tried: %v", defaultConfigsPaths)
}

// displayInstallationResults shows the results of installations to the user
func displayInstallationResults(results map[string]installer.InstallationResult) {
	fmt.Println(resultsHeader)

	successful, failed := displayToolResults(results)
	displaySummary(len(results), successful, failed)
	displayGuidance(successful, failed)
}

// displayToolResults shows individual tool installation results and returns counts
func displayToolResults(results map[string]installer.InstallationResult) (successful, failed int) {
	for toolName, result := range results {
		if result.Success {
			fmt.Printf("%s %s (%s) - installed successfully\n", successIcon, result.Tool.DisplayName, toolName)
			successful++
		} else {
			fmt.Printf("%s %s (%s) - installation failed: %v\n", failureIcon, result.Tool.DisplayName, toolName, result.Error)
			failed++
		}
	}
	return successful, failed
}

// displaySummary shows installation summary statistics
func displaySummary(total, successful, failed int) {
	fmt.Printf(summaryHeader + "\n")
	fmt.Printf("Total attempted: %d\n", total)
	fmt.Printf("Successful: %d\n", successful)
	fmt.Printf("Failed: %d\n", failed)
}

// displayGuidance provides next-step guidance based on installation results
func displayGuidance(successful, failed int) {
	if failed > 0 {
		fmt.Printf("\n" + failureMsg + "\n")
		fmt.Printf("- Run '%s' to check current tool status\n", statusCmdStr)
		fmt.Printf("- Re-run '%s' to retry failed installations\n", retryCmd)
	} else if successful > 0 {
		fmt.Printf("\n" + successMsg + "\n")
		fmt.Printf("Run '%s' to verify your development environment.\n", statusCmdStr)
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
}

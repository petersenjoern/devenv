package cmd

import (
	"fmt"

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
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
}

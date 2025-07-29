package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devenv",
	Short: "DevEnv - Automated developer environment setup",
	Long:  `DevEnv is a Go-based command-line application that automates developer environment setup for personal use. It supports WSL (Ubuntu) and native Linux (Ubuntu) environments.`,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Launch interactive installation of development tools",
	Long:  `Launch interactive TUI for tool selection and installation. First prompts for environment selection (WSL/Linux), then displays categorized tool selection with dependency resolution.`,
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
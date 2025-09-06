package cmd

import (
	"fmt"
	"strings"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/detector"
	"github.com/spf13/cobra"
)

var verbose bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display installation status for all tools",
	Long: `Display table showing installation status for all tools.
Shows binary installation status, configuration status, versions, and paths.
Use --verbose flag for detailed output.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Status command (verbose) - implementation coming soon")
		} else {
			fmt.Println("Status command - implementation coming soon")
		}
	},
}

func GenerateStatusTable(cfg config.Config, det *detector.Detector, verbose bool) string {
	var output strings.Builder
	
	writeHeader(&output, verbose)
	
	for _, categoryTools := range cfg.Categories {
		for _, tool := range categoryTools {
			status := det.DetectTool(tool)
			writeToolRow(&output, tool, status, verbose)
		}
	}
	
	return output.String()
}

func writeHeader(output *strings.Builder, verbose bool) {
	if verbose {
		output.WriteString("Tool Name          Binary    Config    Version        Path\n")
		output.WriteString("-----------------------------------------------------------------\n")
	} else {
		output.WriteString("Tool Name          Binary    Config    Version\n")
		output.WriteString("----------------------------------------------\n")
	}
}

func writeToolRow(output *strings.Builder, tool config.ToolConfig, status detector.Status, verbose bool) {
	binaryStatus := formatStatus(status.BinaryInstalled)
	configStatus := formatStatus(status.ConfigApplied)
	version := formatValue(status.Version)
	
	if verbose {
		path := formatValue(status.Path)
		output.WriteString(fmt.Sprintf("%-18s %-9s %-9s %-14s %s\n",
			tool.DisplayName, binaryStatus, configStatus, version, path))
	} else {
		output.WriteString(fmt.Sprintf("%-18s %-9s %-9s %s\n",
			tool.DisplayName, binaryStatus, configStatus, version))
	}
}

func formatStatus(installed bool) string {
	if installed {
		return "✓"
	}
	return "✗"
}

func formatValue(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func init() {
	statusCmd.Flags().BoolVar(&verbose, "verbose", false,
		"Verbose output including installation paths and version information")
	rootCmd.AddCommand(statusCmd)
}

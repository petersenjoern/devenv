package cmd

import (
	"fmt"
	"strings"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/detector"
	"github.com/spf13/cobra"
)

const (
	normalHeader     = "Tool Name          Binary    Config    Version\n"
	normalSeparator  = "----------------------------------------------\n"
	verboseHeader    = "Tool Name          Binary    Config    Version        Path\n"
	verboseSeparator = "-----------------------------------------------------------------\n"
	toolNameWidth    = 18
	statusWidth      = 9
	versionWidth     = 14
)

var verbose bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display installation status for all tools",
	Long: `Display table showing installation status for all tools.
Shows binary installation status, configuration status, versions, and paths.
Use --verbose flag for detailed output.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := executeStatusCommand(verbose)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func executeStatusCommand(verbose bool) error {
	configPath, err := findConfigPath()
	if err != nil {
		return fmt.Errorf("finding config: %w", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	detector := detector.New()
	statusTable := GenerateStatusTable(cfg, detector, verbose)
	fmt.Print(statusTable)
	return nil
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
		output.WriteString(verboseHeader)
		output.WriteString(verboseSeparator)
	} else {
		output.WriteString(normalHeader)
		output.WriteString(normalSeparator)
	}
}

func writeToolRow(output *strings.Builder, tool config.ToolConfig, status detector.Status, verbose bool) {
	binaryStatus := formatStatus(status.BinaryInstalled)
	configStatus := formatStatus(status.ConfigApplied)
	version := formatValue(status.Version)

	if verbose {
		path := formatValue(status.Path)
		output.WriteString(fmt.Sprintf("%-*s %-*s %-*s %-*s %s\n",
			toolNameWidth, tool.DisplayName,
			statusWidth, binaryStatus,
			statusWidth, configStatus,
			versionWidth, version,
			path))
	} else {
		output.WriteString(fmt.Sprintf("%-*s %-*s %-*s %s\n",
			toolNameWidth, tool.DisplayName,
			statusWidth, binaryStatus,
			statusWidth, configStatus,
			version))
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

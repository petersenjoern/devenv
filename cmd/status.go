package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verbose bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display installation status for all tools",
	Long:  `Display table showing installation status for all tools. Shows binary installation status, configuration status, versions, and paths. Use --verbose flag for detailed output.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Status command (verbose) - implementation coming soon")
		} else {
			fmt.Println("Status command - implementation coming soon")
		}
	},
}

func init() {
	statusCmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose output including installation paths and version information")
	rootCmd.AddCommand(statusCmd)
}
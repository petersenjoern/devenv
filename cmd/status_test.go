package cmd

import (
	"strings"
	"testing"

	"github.com/petersenjoern/devenv/internal/config"
	"github.com/petersenjoern/devenv/internal/detector"
)

func TestStatusCommand_ShouldDisplayTableWithToolStatus(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
					ConfigPath:  "/etc/bash.bashrc",
				},
			},
		},
	}

	detector := detector.New()

	output := GenerateStatusTable(cfg, detector, false)

	if output == "" {
		t.Errorf("Expected status table output, got empty string")
	}

	if !strings.Contains(output, "Tool Name") {
		t.Errorf("Expected status table to contain header 'Tool Name', got: %s", output)
	}

	if !strings.Contains(output, "Bash Shell") {
		t.Errorf("Expected status table to contain 'Bash Shell', got: %s", output)
	}
}

func TestStatusCommand_ShouldDisplayVerboseOutput(t *testing.T) {
	cfg := config.Config{
		Categories: map[string]config.CategoryConfig{
			"shells": {
				"bash": config.ToolConfig{
					DisplayName: "Bash Shell",
					BinaryName:  "bash",
					ConfigPath:  "/etc/bash.bashrc",
				},
			},
		},
	}

	detector := detector.New()

	normalOutput := GenerateStatusTable(cfg, detector, false)
	verboseOutput := GenerateStatusTable(cfg, detector, true)

	if len(verboseOutput) <= len(normalOutput) {
		t.Errorf("Expected verbose output to be longer than normal output")
	}
}

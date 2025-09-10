package detector

import (
	"os"
	"os/exec"
	"strings"

	"github.com/petersenjoern/devenv/internal/config"
)

type Status struct {
	BinaryInstalled bool
	ConfigApplied   bool
	Version         string
	Path            string
}

type Detector struct{}

func New() *Detector {
	return &Detector{}
}

func (d *Detector) DetectTool(tool config.ToolConfig) Status {
	path, err := exec.LookPath(tool.BinaryName)

	if err != nil {
		return Status{
			BinaryInstalled: false,
			ConfigApplied:   false,
			Version:         "",
			Path:            "",
		}
	}

	return Status{
		BinaryInstalled: true,
		ConfigApplied:   d.IsConfigExisting(tool.ConfigPath),
		Version:         d.GetVersion(tool.BinaryName),
		Path:            path,
	}
}

func (d *Detector) DetectEnvironment() (string, error) {
	return "linux", nil
}

func (d *Detector) IsBinaryInstalled(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

func (d *Detector) GetVersion(binaryName string) string {
	cmd := exec.Command(binaryName, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	version := strings.TrimSpace(string(output))
	lines := strings.Split(version, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return "unknown"
}

func (d *Detector) IsConfigExisting(configPath string) bool {
	if configPath == "" {
		return false
	}

	_, err := os.Stat(configPath)
	return err == nil
}

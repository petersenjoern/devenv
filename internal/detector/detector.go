package detector

import (
	"os/exec"

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
	return Status{
		BinaryInstalled: d.IsBinaryInstalled(tool.BinaryName),
		ConfigApplied:   false,
		Version:         "",
		Path:            "",
	}
}

func (d *Detector) DetectEnvironment() (string, error) {
	return "linux", nil
}

func (d *Detector) IsBinaryInstalled(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

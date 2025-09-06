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
	var binaryPath string
	binaryInstalled := d.IsBinaryInstalled(tool.BinaryName)

	if binaryInstalled {
		path, _ := exec.LookPath(tool.BinaryName)
		binaryPath = path
	} else {
		binaryPath = ""
	}

	return Status{
		BinaryInstalled: binaryInstalled,
		ConfigApplied:   false,
		Version:         "",
		Path:            binaryPath,
	}
}

func (d *Detector) DetectEnvironment() (string, error) {
	return "linux", nil
}

func (d *Detector) IsBinaryInstalled(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

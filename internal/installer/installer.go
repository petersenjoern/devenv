package installer

import (
	"github.com/petersenjoern/devenv/internal/config"
)

type Installer interface {
	Install(tool config.ToolConfig) error
}

type APTInstaller struct{}
type ScriptInstaller struct{}
type ManualInstaller struct{}

func (a *APTInstaller) Install(tool config.ToolConfig) error {
	return nil
}

func (s *ScriptInstaller) Install(tool config.ToolConfig) error {
	return nil
}

func (m *ManualInstaller) Install(tool config.ToolConfig) error {
	return nil
}

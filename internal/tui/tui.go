package tui

import (
	"github.com/jrjl/devenv/internal/config"
)

type TUI struct {
	config config.Config
}

func New(cfg config.Config) *TUI {
	return &TUI{
		config: cfg,
	}
}

func (t *TUI) ShowEnvironmentSelection() (string, error) {
	return "linux", nil
}

func (t *TUI) ShowToolSelection() ([]string, error) {
	return []string{}, nil
}

func (t *TUI) ShowInstallationProgress(tools []string) error {
	return nil
}
package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/petersenjoern/devenv/internal/config"
)

type TUI struct {
	config config.Config
}

type Selections struct {
	CategoryAndTools []CategoryAndTools
}

type CategoryAndTools struct {
	Category string
	Tools    []string
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

func CreateInstallationForm() (Selections, error) {
	var selections Selections

	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		return selections, fmt.Errorf("failed to load config: %w", err)
	}
	categories := config.GetCategories(cfg)
	if len(categories) == 0 {
		return selections, fmt.Errorf("no categories found in config")
	}

	var huhGroups []*huh.Group
	for _, category := range categories {
		selection := CategoryAndTools{
			Category: category,
			Tools:    []string{},
		}
		tools, exists := config.GetToolsInCategory(cfg, category)
		if !exists {
			return selections, fmt.Errorf("category %s does not exist", category)
		}
		for toolName := range tools {
			tool, exists := config.GetTool(cfg, category, toolName)
			if !exists {
				return selections, fmt.Errorf("tool %s does not exist in category %s", toolName, category)
			}
			selection.Tools = append(selection.Tools, tool.BinaryName)
		}

		var options []huh.Option[string]
		for toolName := range tools {
			tool, exists := config.GetTool(cfg, category, toolName)
			if !exists {
				return selections, fmt.Errorf("tool %s does not exist in category %s", toolName, category)
			}
			options = append(options, huh.NewOption(tool.DisplayName, tool.BinaryName))
		}
		huhGroupPerCategory := huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(category).
				Options(options...).
				Value(&selection.Tools),
		)
		huhGroups = append(huhGroups, huhGroupPerCategory)
		selections.CategoryAndTools = append(selections.CategoryAndTools, selection)
	}

	form := huh.NewForm(huhGroups...)
	err = form.Run()
	if err != nil {
		return selections, fmt.Errorf("error running installation form: %w", err)
	}

	return selections, nil
}

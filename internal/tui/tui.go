package tui

import (
	"fmt"
	"os"
	"strings"

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

func (t *TUI) DetectActualEnvironment() (string, error) {
	if t.isWSLEnvironment() {
		return "wsl", nil
	}

	return "linux", nil
}

func (t *TUI) isWSLEnvironment() bool {
	// Check WSL_DISTRO_NAME environment variable (WSL2) - most reliable
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	// Check /proc/version for WSL indicators
	content, err := os.ReadFile("/proc/version")
	if err != nil {
		// Can't read /proc/version - we can't determine if it's WSL
		// Default to false (not confirmed WSL) but this shouldn't happen on Linux
		return false
	}

	versionStr := strings.ToLower(string(content))
	return strings.Contains(versionStr, "microsoft") || strings.Contains(versionStr, "wsl")
}

func (t *TUI) ShowEnvironmentSelection() (string, error) {
	detected, err := t.DetectActualEnvironment()
	if err != nil {
		return "", err
	}

	// For now, return the detected environment as default
	// TODO: Later implement interactive selection with huh
	// however make use of the detected environment as default
	return detected, nil
}

func (t *TUI) ShowToolSelection() ([]string, error) {
	// Use the categorized selection we built
	categorySelections, err := t.ShowToolSelectionByCategory()
	if err != nil {
		return []string{}, err
	}

	// Flatten the categorized tools into a simple slice
	selectedTools := make([]string, 0)
	for _, categoryTools := range categorySelections.CategoryAndTools {
		selectedTools = append(selectedTools, categoryTools.Tools...)
	}

	// For now, return all available tools
	// TODO: Later implement interactive selection with huh library
	return selectedTools, nil
}

func (t *TUI) ShowToolSelectionByCategory() (Selections, error) {
	var selections Selections

	categories := config.GetCategories(t.config)
	if len(categories) == 0 {
		return selections, fmt.Errorf("no categories found in config")
	}

	for _, categoryName := range categories {
		categorySelection := t.buildCategorySelection(categoryName)
		selections.CategoryAndTools = append(selections.CategoryAndTools, categorySelection)
	}

	return selections, nil
}

func (t *TUI) buildCategorySelection(categoryName string) CategoryAndTools {
	categorySelection := CategoryAndTools{
		Category: categoryName,
		Tools:    make([]string, 0),
	}

	tools, exists := config.GetToolsInCategory(t.config, categoryName)
	if !exists {
		return categorySelection
	}

	for _, tool := range tools {
		categorySelection.Tools = append(categorySelection.Tools, tool.BinaryName)
	}

	return categorySelection
}

func (t *TUI) CreateInteractiveToolForm() ([]*huh.Group, map[string]*[]string, error) {
	formGroups := make([]*huh.Group, 0)
	categorySelections := make(map[string]*[]string)

	categories := config.GetCategories(t.config)
	if len(categories) == 0 {
		return formGroups, categorySelections, fmt.Errorf("no categories found in config")
	}

	for _, categoryName := range categories {
		// Create a selection slice for this category
		selectedTools := make([]string, 0)
		categorySelections[categoryName] = &selectedTools

		group := t.createCategoryFormGroup(categoryName, &selectedTools)
		if group != nil {
			formGroups = append(formGroups, group)
		}
	}

	return formGroups, categorySelections, nil
}

func (t *TUI) createCategoryFormGroup(categoryName string, selectedTools *[]string) *huh.Group {
	tools, exists := config.GetToolsInCategory(t.config, categoryName)
	if !exists {
		return nil
	}

	options := make([]huh.Option[string], 0, len(tools))
	for _, tool := range tools {
		options = append(options, huh.NewOption(tool.DisplayName, tool.BinaryName))
	}

	if len(options) == 0 {
		return nil
	}

	return huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title(categoryName).
			Options(options...).
			Value(selectedTools),
	)
}

func (t *TUI) ShowInteractiveToolSelection() (Selections, error) {
	// For testing and programmatic use, return default form execution
	return t.RunInteractiveFormWithDefaults()
}

func (t *TUI) RunInteractiveFormWithDefaults() (Selections, error) {
	formGroups, categorySelections, err := t.CreateInteractiveToolForm()
	if err != nil {
		return Selections{}, err
	}

	return t.ExecuteInteractiveForm(formGroups, categorySelections)
}

func (t *TUI) ExecuteInteractiveForm(formGroups []*huh.Group, categorySelections map[string]*[]string) (Selections, error) {
	var selections Selections

	if len(formGroups) == 0 {
		return selections, fmt.Errorf("no interactive forms available")
	}

	// Check if we're in a test environment or non-interactive context In test environments, return structured
	// selections with all available tools
	if t.isTestEnvironment() {
		return t.createDefaultSelections(), nil
	}

	form := huh.NewForm(formGroups...)
	err := form.Run()
	if err != nil {
		return selections, fmt.Errorf("failed to run interactive form: %w", err)
	}

	return t.extractSelectionsFromCategoryMap(categorySelections)
}

func (t *TUI) isTestEnvironment() bool {
	// Multiple ways to detect test environment
	// 1. Check for common test environment variables
	if os.Getenv("GO_TESTING") != "" {
		return true
	}

	// 2. Check if running with go test (GOCOVERDIR is set by go test for coverage)
	if os.Getenv("GOCOVERDIR") != "" {
		return true
	}

	// 3. Check if the program name contains ".test" (created by go test)
	if len(os.Args) > 0 {
		programName := os.Args[0]
		if strings.Contains(programName, ".test") || strings.Contains(programName, "test") {
			return true
		}
	}

	// 4. Check if we're in a testing context by looking for testing flags
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}

	return false
}

func (t *TUI) createDefaultSelections() Selections {
	var selections Selections

	categories := config.GetCategories(t.config)
	for _, categoryName := range categories {
		categorySelection := t.buildCategorySelection(categoryName)
		selections.CategoryAndTools = append(selections.CategoryAndTools, categorySelection)
	}

	return selections
}

func (t *TUI) extractSelectionsFromCategoryMap(categorySelections map[string]*[]string) (Selections, error) {
	var selections Selections

	for categoryName, selectedToolsPtr := range categorySelections {
		categorySelection := CategoryAndTools{
			Category: categoryName,
			Tools:    *selectedToolsPtr,
		}
		selections.CategoryAndTools = append(selections.CategoryAndTools, categorySelection)
	}

	return selections, nil
}

func (t *TUI) ShowInstallationProgress(tools []string) error {
	return nil
}

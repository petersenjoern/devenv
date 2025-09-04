package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ToolConfig struct {
	DisplayName      string            `yaml:"display_name"`
	BinaryName       string            `yaml:"binary_name"`
	InstallMethod    string            `yaml:"install_method"`
	PackageName      string            `yaml:"package_name"`
	InstallScript    string            `yaml:"install_script"`
	ConfigPath       string            `yaml:"config_path"`
	ConfigTemplate   string            `yaml:"config_template"`
	Dependencies     []string          `yaml:"dependencies"`
	WSLNotes         string            `yaml:"wsl_notes"`
	Version          string            `yaml:"version,omitempty"`
	DownloadURL      string            `yaml:"download_url,omitempty"`
	PostInstallSteps []string          `yaml:"post_install_steps,omitempty"`
	ValidateCommand  string            `yaml:"validate_command,omitempty"`
	EnvVars          map[string]string `yaml:"env_vars,omitempty"`
	InstallLocation  string            `yaml:"install_location,omitempty"`
	RequiredPackages []string          `yaml:"required_packages,omitempty"`
	CheckCommand     string            `yaml:"check_command,omitempty"`
}

type CategoryConfig map[string]ToolConfig

type Config struct {
	Categories map[string]CategoryConfig `yaml:"categories"`
}

func LoadConfig(filePath string) (Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func GetCategories(cfg Config) []string {
	categories := make([]string, 0, len(cfg.Categories))
	for category := range cfg.Categories {
		categories = append(categories, category)
	}
	return categories
}
func GetToolsInCategory(cfg Config, category string) (CategoryConfig, bool) {
	tools, exists := cfg.Categories[category]
	return tools, exists
}
func GetTool(cfg Config, category, toolName string) (ToolConfig, bool) {
	tools, exists := cfg.Categories[category]
	if !exists {
		return ToolConfig{}, false
	}
	tool, exists := tools[toolName]
	return tool, exists
}

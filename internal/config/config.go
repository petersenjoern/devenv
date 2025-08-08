package config

import (
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
		return Config{}, err
	}
	
	return cfg, nil
}

func loadConfig() (Config, error) {
	return LoadConfig("config.yaml")
}

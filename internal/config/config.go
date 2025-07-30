package config

type ToolConfig struct {
	DisplayName    string   `yaml:"display_name"`
	BinaryName     string   `yaml:"binary_name"`
	InstallMethod  string   `yaml:"install_method"`
	PackageName    string   `yaml:"package_name"`
	InstallScript  string   `yaml:"install_script"`
	ConfigPath     string   `yaml:"config_path"`
	ConfigTemplate string   `yaml:"config_template"`
	Dependencies   []string `yaml:"dependencies"`
	WSLNotes       string   `yaml:"wsl_notes"`
}

type CategoryConfig map[string]ToolConfig

type Config struct {
	Categories map[string]CategoryConfig `yaml:"categories"`
}

package config

import (
	"testing"
)

func TestLoadConfig_ShouldLoadValidYAMLFile(t *testing.T) {
	config, err := LoadConfig("../../config.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Categories == nil {
		t.Fatal("Expected categories to be loaded, got nil")
	}

	if len(config.Categories) == 0 {
		t.Fatal("Expected categories to contain tools, got empty map")
	}

	terminals, exists := config.Categories["terminals"]
	if !exists {
		t.Fatal("Expected 'terminals' category to exist")
	}

	alacritty, exists := terminals["alacritty"]
	if !exists {
		t.Fatal("Expected 'alacritty' tool to exist in terminals category")
	}

	if alacritty.DisplayName != "Alacritty Terminal" {
		t.Errorf("Expected display_name 'Alacritty Terminal', got '%s'", alacritty.DisplayName)
	}

	if alacritty.BinaryName != "alacritty" {
		t.Errorf("Expected binary_name 'alacritty', got '%s'", alacritty.BinaryName)
	}

	if alacritty.InstallMethod != "manual" {
		t.Errorf("Expected install_method 'manual', got '%s'", alacritty.InstallMethod)
	}
}

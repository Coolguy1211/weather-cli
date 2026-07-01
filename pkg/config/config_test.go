package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadDefault(t *testing.T) {
	// Set temp home directory
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	// Since the config file won't exist in temp home, Load should return defaults
	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error loading defaults, got: %v", err)
	}

	if cfg.DefaultCity != "" {
		t.Errorf("expected empty DefaultCity, got: %q", cfg.DefaultCity)
	}

	if cfg.Unit != "celsius" {
		t.Errorf("expected default unit 'celsius', got: %q", cfg.Unit)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Set temp home directory
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	expectedCity := "San Francisco"
	expectedUnit := "fahrenheit"

	cfg := &Config{
		DefaultCity: expectedCity,
		Unit:        expectedUnit,
	}

	err := Save(cfg)
	if err != nil {
		t.Fatalf("expected no error saving config, got: %v", err)
	}

	// Verify file was created
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created at expected path: %s", configPath)
	}

	// Load and verify values
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error loading config, got: %v", err)
	}

	if loadedCfg.DefaultCity != expectedCity {
		t.Errorf("expected DefaultCity %q, got %q", expectedCity, loadedCfg.DefaultCity)
	}

	if loadedCfg.Unit != expectedUnit {
		t.Errorf("expected Unit %q, got %q", expectedUnit, loadedCfg.Unit)
	}
}

func TestConfigSaveInvalidUnit(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	cfg := &Config{
		DefaultCity: "London",
		Unit:        "invalid_unit",
	}

	err := Save(cfg)
	if err == nil {
		t.Error("expected error when saving invalid unit, got nil")
	}
}

func TestConfigPath(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("expected no error getting config path, got: %v", err)
	}

	expectedPath := filepath.Join(tempHome, ".weather-config.json")
	if path != expectedPath {
		t.Errorf("expected config path %q, got %q", expectedPath, path)
	}
}

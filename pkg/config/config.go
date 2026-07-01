package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DefaultCity string `json:"default_city"`
	Unit        string `json:"unit"` // "celsius", "fahrenheit", "kelvin"
}

// GetConfigPath returns the path to the ~/.weather-config.json file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}
	return filepath.Join(home, ".weather-config.json"), nil
}

// Load reads the configuration from the user's home directory.
// If the file does not exist, it returns a default configuration.
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			DefaultCity: "",
			Unit:        "celsius",
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Normalize
	cfg.Unit = strings.ToLower(cfg.Unit)
	if cfg.Unit == "" {
		cfg.Unit = "celsius"
	}

	return &cfg, nil
}

// Save writes the configuration back to the ~/.weather-config.json file.
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Normalize unit
	cfg.Unit = strings.ToLower(cfg.Unit)
	if cfg.Unit != "celsius" && cfg.Unit != "fahrenheit" && cfg.Unit != "kelvin" {
		return fmt.Errorf("invalid unit %q, must be 'celsius', 'fahrenheit', or 'kelvin'", cfg.Unit)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

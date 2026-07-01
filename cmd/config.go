package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/friendly-maxwell/weather/pkg/config"
)

var (
	configCity string
	configUnit string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or modify configuration settings",
	Long:  `View current configuration values or update settings like the default city and temperature unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		path, err := config.GetConfigPath()
		if err != nil {
			return err
		}

		// If no modifications are requested, show current settings
		if configCity == "" && configUnit == "" {
			fmt.Printf("📂 Configuration File: %s\n", path)
			fmt.Printf("🌆 Default City:       %s\n", getDisplayVal(cfg.DefaultCity))
			fmt.Printf("🌡️  Temperature Unit:  %s\n", cfg.Unit)
			return nil
		}

		// Apply modifications
		modified := false
		if configCity != "" {
			cfg.DefaultCity = configCity
			modified = true
			fmt.Printf("🌆 Default city set to %q\n", configCity)
		}

		if configUnit != "" {
			unitLower := strings.ToLower(configUnit)
			if unitLower != "celsius" && unitLower != "fahrenheit" && unitLower != "kelvin" {
				return fmt.Errorf("invalid unit %q; must be 'celsius', 'fahrenheit', or 'kelvin'", configUnit)
			}
			cfg.Unit = unitLower
			modified = true
			fmt.Printf("🌡️  Default unit set to %q\n", unitLower)
		}

		if modified {
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}
			fmt.Println("💾 Configuration saved successfully.")
		}

		return nil
	},
}

func getDisplayVal(val string) string {
	if val == "" {
		return "(not set, will use IP geolocation)"
	}
	return val
}

func init() {
	configCmd.Flags().StringVarP(&configCity, "city", "c", "", "Set default city")
	configCmd.Flags().StringVarP(&configUnit, "unit", "u", "", "Set default temperature unit (celsius, fahrenheit, kelvin)")
	rootCmd.AddCommand(configCmd)
}

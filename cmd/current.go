package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/friendly-maxwell/weather/pkg/config"
	"github.com/friendly-maxwell/weather/pkg/formatter"
	"github.com/friendly-maxwell/weather/pkg/weather"
)

var cityFlagCurrent string

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current weather conditions",
	Long:  `Fetches and displays the current temperature, humidity, wind speed, and weather condition.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		unit := cfg.Unit
		if unitFlag != "" {
			unit = unitFlag
		}

		var lat, lon float64
		var resolvedCityName string

		// 1. Resolve coordinates
		if cityFlagCurrent != "" {
			fmt.Printf("🔍 Resolving location coordinates for %q...\n", cityFlagCurrent)
			geo, err := weather.ResolveCity(cityFlagCurrent)
			if err != nil {
				return err
			}
			lat = geo.Latitude
			lon = geo.Longitude
			region := geo.Admin1
			if region == "" {
				resolvedCityName = fmt.Sprintf("%s, %s", geo.Name, geo.Country)
			} else {
				resolvedCityName = fmt.Sprintf("%s, %s (%s)", geo.Name, geo.Country, region)
			}
		} else if cfg.DefaultCity != "" {
			fmt.Printf("🔍 Resolving coordinates for default city %q...\n", cfg.DefaultCity)
			geo, err := weather.ResolveCity(cfg.DefaultCity)
			if err != nil {
				return err
			}
			lat = geo.Latitude
			lon = geo.Longitude
			region := geo.Admin1
			if region == "" {
				resolvedCityName = fmt.Sprintf("%s, %s", geo.Name, geo.Country)
			} else {
				resolvedCityName = fmt.Sprintf("%s, %s (%s)", geo.Name, geo.Country, region)
			}
		} else {
			fmt.Println("📍 No city specified. Detecting location by IP address...")
			loc, err := weather.DetectLocation()
			if err != nil {
				return fmt.Errorf("failed to auto-detect location (please specify a city with --city): %w", err)
			}
			lat = loc.Latitude
			lon = loc.Longitude
			resolvedCityName = fmt.Sprintf("%s (Auto-detected)", loc.City)
		}

		// 2. Fetch weather
		fmt.Printf("🌤️  Fetching weather for %s...\n", resolvedCityName)
		wResp, err := weather.GetWeather(lat, lon, unit)
		if err != nil {
			return err
		}

		// 3. Print current weather
		formatter.PrintCurrentWeather(resolvedCityName, wResp, unit)
		return nil
	},
}

func init() {
	currentCmd.Flags().StringVarP(&cityFlagCurrent, "city", "c", "", "City name to check weather for")
	rootCmd.AddCommand(currentCmd)
}

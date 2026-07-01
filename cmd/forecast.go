package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/friendly-maxwell/weather/pkg/config"
	"github.com/friendly-maxwell/weather/pkg/formatter"
	"github.com/friendly-maxwell/weather/pkg/weather"
)

var (
	cityFlagForecast string
	daysFlag         int
)

var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Show weather forecast for the next few days",
	Long:  `Fetches and displays the daily weather forecast, including conditions, minimum, and maximum temperatures.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if daysFlag < 1 || daysFlag > 5 {
			return fmt.Errorf("invalid number of days: %d (must be between 1 and 5)", daysFlag)
		}

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
		if cityFlagForecast != "" {
			fmt.Printf("🔍 Resolving location coordinates for %q...\n", cityFlagForecast)
			geo, err := weather.ResolveCity(cityFlagForecast)
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
		fmt.Printf("🌤️  Fetching daily forecast for %s...\n", resolvedCityName)
		wResp, err := weather.GetWeather(lat, lon, unit)
		if err != nil {
			return err
		}

		// 3. Print forecast
		formatter.PrintForecast(resolvedCityName, wResp, unit, daysFlag)
		return nil
	},
}

func init() {
	forecastCmd.Flags().StringVarP(&cityFlagForecast, "city", "c", "", "City name to check weather for")
	forecastCmd.Flags().IntVarP(&daysFlag, "days", "d", 3, "Number of days for forecast (1-5)")
	rootCmd.AddCommand(forecastCmd)
}

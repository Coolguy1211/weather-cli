package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/friendly-maxwell/weather/pkg/weather"
)

type WeatherStyle struct {
	Description string
	Icon        string
	ColorAttr   *color.Color
}

// GetWeatherStyle maps a WMO weather code to text description, icon, and terminal color
func GetWeatherStyle(code int) WeatherStyle {
	var desc, icon string
	var colorAttr *color.Color

	yellow := color.New(color.FgYellow, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)
	white := color.New(color.FgWhite, color.Bold)
	gray := color.New(color.FgHiBlack)
	red := color.New(color.FgRed, color.Bold)

	switch code {
	case 0:
		desc, icon, colorAttr = "Clear sky", "☀️", yellow
	case 1:
		desc, icon, colorAttr = "Mainly clear", "🌤️", yellow
	case 2:
		desc, icon, colorAttr = "Partly cloudy", "⛅", cyan
	case 3:
		desc, icon, colorAttr = "Overcast", "☁️", white
	case 45, 48:
		desc, icon, colorAttr = "Fog", "🌫️", gray
	case 51, 53, 55:
		desc, icon, colorAttr = "Drizzle", "🌧️", blue
	case 56, 57:
		desc, icon, colorAttr = "Freezing Drizzle", "🌧️", blue
	case 61, 63, 65:
		desc, icon, colorAttr = "Rain", "🌧️", blue
	case 66, 67:
		desc, icon, colorAttr = "Freezing Rain", "🌧️", blue
	case 71, 73, 75:
		desc, icon, colorAttr = "Snow fall", "❄️", cyan
	case 77:
		desc, icon, colorAttr = "Snow grains", "❄️", cyan
	case 80, 81, 82:
		desc, icon, colorAttr = "Rain showers", "🌧️", blue
	case 85, 86:
		desc, icon, colorAttr = "Snow showers", "❄️", cyan
	case 95:
		desc, icon, colorAttr = "Thunderstorm", "⛈️", red
	case 96, 99:
		desc, icon, colorAttr = "Thunderstorm with hail", "⛈️", red
	default:
		desc, icon, colorAttr = "Unknown weather", "❓", white
	}

	return WeatherStyle{
		Description: desc,
		Icon:        icon,
		ColorAttr:   colorAttr,
	}
}

// formatTemp converts and formats the temperature according to requested unit
func formatTemp(celsiusVal float64, apiVal float64, unit string) string {
	switch strings.ToLower(unit) {
	case "kelvin":
		return fmt.Sprintf("%.1f K", celsiusVal+273.15)
	case "fahrenheit":
		// If apiVal is in Fahrenheit (pre-converted by API), use it directly.
		// If not, we convert it ourselves.
		// But in our client.go, if unit == "fahrenheit", the API response temperature is already in Fahrenheit.
		// So we just use apiVal.
		return fmt.Sprintf("%.1f°F", apiVal)
	default:
		return fmt.Sprintf("%.1f°C", apiVal)
	}
}

// PrintCurrentWeather prints current weather details in a beautiful terminal layout
func PrintCurrentWeather(locationName string, w *weather.WeatherResponse, unit string) {
	style := GetWeatherStyle(w.Current.WeatherCode)

	cVal := w.Current.Temperature2m
	appCVal := w.Current.ApparentTemperature
	// If the API returned Fahrenheit, we need Celsius for Kelvin calculations.
	// Otherwise, cVal is already Celsius.
	if w.CurrentUnits.Temperature2m == "°F" {
		cVal = (w.Current.Temperature2m - 32) * 5 / 9
		appCVal = (w.Current.ApparentTemperature - 32) * 5 / 9
	}

	tempStr := formatTemp(cVal, w.Current.Temperature2m, unit)
	apparentStr := formatTemp(appCVal, w.Current.ApparentTemperature, unit)

	headerColor := color.New(color.FgMagenta, color.Bold)
	labelColor := color.New(color.FgCyan)
	valueColor := color.New(color.FgWhite)

	separator := strings.Repeat("═", 50)
	fmt.Println(headerColor.Sprint(separator))
	fmt.Printf("%s %s\n", headerColor.Sprint("  Weather Report:"), color.New(color.Bold, color.FgHiWhite).Sprint(locationName))
	fmt.Println(headerColor.Sprint(separator))

	// Display weather with colored icon
	fmt.Printf("  %-18s: ", labelColor.Sprint("Condition"))
	style.ColorAttr.Printf("%s  %s\n", style.Icon, style.Description)

	fmt.Printf("  %-18s: %s\n", labelColor.Sprint("Temperature"), valueColor.Sprintf("%s (Feels like %s)", tempStr, apparentStr))
	fmt.Printf("  %-18s: %s%%\n", labelColor.Sprint("Humidity"), valueColor.Sprintf("%d", w.Current.RelativeHumidity2m))
	fmt.Printf("  %-18s: %s %s\n", labelColor.Sprint("Wind Speed"), valueColor.Sprintf("%.1f", w.Current.WindSpeed10m), w.CurrentUnits.WindSpeed10m)
	fmt.Println(headerColor.Sprint(separator))
}

// PrintForecast prints daily forecast details in a formatted table
func PrintForecast(locationName string, w *weather.WeatherResponse, unit string, numDays int) {
	if numDays <= 0 || numDays > len(w.Daily.Time) {
		numDays = 3
	}

	headerColor := color.New(color.FgBlue, color.Bold)
	tableHeaderColor := color.New(color.FgHiCyan, color.Underline)
	separatorColor := color.New(color.FgHiBlack)

	fmt.Println()
	separator := strings.Repeat("─", 65)
	fmt.Println(headerColor.Sprint(separator))
	fmt.Printf("%s %d-Day Forecast for %s\n", headerColor.Sprint(" 📅"), numDays, color.New(color.Bold, color.FgWhite).Sprint(locationName))
	fmt.Println(headerColor.Sprint(separator))

	// Print headers
	fmt.Printf("  %-12s  %-24s  %-10s  %-10s\n",
		tableHeaderColor.Sprint("Date"),
		tableHeaderColor.Sprint("Condition"),
		tableHeaderColor.Sprint("Min Temp"),
		tableHeaderColor.Sprint("Max Temp"),
	)
	fmt.Println(separatorColor.Sprint(separator))

	for i := 0; i < numDays; i++ {
		dateStr := w.Daily.Time[i]
		t, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			dateStr = t.Format("Mon, Jan 02")
		}

		style := GetWeatherStyle(w.Daily.WeatherCode[i])

		// Base values for Kelvin conversion
		minVal := w.Daily.Temperature2mMin[i]
		maxVal := w.Daily.Temperature2mMax[i]
		cMinVal := minVal
		cMaxVal := maxVal
		if w.CurrentUnits.Temperature2m == "°F" {
			cMinVal = (minVal - 32) * 5 / 9
			cMaxVal = (maxVal - 32) * 5 / 9
		}

		minStr := formatTemp(cMinVal, minVal, unit)
		maxStr := formatTemp(cMaxVal, maxVal, unit)

		// Print Row
		fmt.Printf("  %-12s  %s %-20s  %-10s  %-10s\n",
			color.WhiteString(dateStr),
			style.ColorAttr.Sprint(style.Icon),
			style.ColorAttr.Sprint(style.Description),
			color.BlueString(minStr),
			color.RedString(maxStr),
		)
	}
	fmt.Println(headerColor.Sprint(separator))
}

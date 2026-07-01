package formatter

import (
	"testing"
)

func TestGetWeatherStyle(t *testing.T) {
	tests := []struct {
		code         int
		expectedDesc string
		expectedIcon string
	}{
		{0, "Clear sky", "☀️"},
		{1, "Mainly clear", "🌤️"},
		{3, "Overcast", "☁️"},
		{61, "Rain", "🌧️"},
		{71, "Snow fall", "❄️"},
		{95, "Thunderstorm", "⛈️"},
		{999, "Unknown weather", "❓"},
	}

	for _, tc := range tests {
		style := GetWeatherStyle(tc.code)
		if style.Description != tc.expectedDesc {
			t.Errorf("for code %d: expected description %q, got %q", tc.code, tc.expectedDesc, style.Description)
		}
		if style.Icon != tc.expectedIcon {
			t.Errorf("for code %d: expected icon %q, got %q", tc.code, tc.expectedIcon, style.Icon)
		}
	}
}

func TestFormatTemp(t *testing.T) {
	tests := []struct {
		celsiusVal float64
		apiVal     float64
		unit       string
		expected   string
	}{
		{20.0, 20.0, "celsius", "20.0°C"},
		{20.0, 20.0, "CELSIUS", "20.0°C"},
		{20.0, 68.0, "fahrenheit", "68.0°F"}, // API returned Fahrenheit directly
		{20.0, 20.0, "kelvin", "293.1 K"},     // We convert 20.0 Celsius to Kelvin (20 + 273.15 = 293.15 -> rounds to 293.1 K)
	}

	for _, tc := range tests {
		res := formatTemp(tc.celsiusVal, tc.apiVal, tc.unit)
		if res != tc.expected {
			t.Errorf("for unit %q: expected %q, got %q", tc.unit, tc.expected, res)
		}
	}
}

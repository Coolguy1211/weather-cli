# Weather CLI

A fast, clean, and beautiful command-line weather client written in Go. It fetches current conditions and forecasts directly from the free **Open-Meteo API** (no API key required) and utilizes **ip-api.com** for automatic location detection.

---

## Features

- 🌤️ **Current Weather**: Get temperature, apparent temperature ("feels like"), humidity, and wind speed.
- 📅 **Multi-Day Forecast**: View a detailed 1 to 5-day forecast rendered in a clean terminal table.
- 🎨 **Rich Output**: Colored terminal rendering featuring weather condition icons mapped from standard WMO interpretation codes.
- ⚙️ **Persistent Configuration**: Save your default city and temperature unit in `~/.weather-config.json` so you don't have to specify flags every time.
- 📐 **Unit Conversion**: Easily switch between Celsius, Fahrenheit, and Kelvin.
- 📍 **IP Geolocation**: Automatically detects your city based on your public IP address when no city is specified.

---

## Installation

### Prerequisites
Make sure you have **Go** installed on your system. 

### Install via Go Toolchain
Clone the repository and run:
```bash
go install
```
This compiles the code and places the binary in your `GOPATH/bin` folder, making it available globally (assuming `GOPATH/bin` is in your environment `PATH`).

---

## Usage

### 1. Check Today's Weather
Fetch current conditions for a city:
```bash
weather current --city "New York"
```
Or check the weather at your current location (auto-detected via IP):
```bash
weather current
```

### 2. Check the Forecast
View a 3-day forecast for a city:
```bash
weather forecast --city "Tokyo" --days 3
```

### 3. Manage Settings
Save your preferred city and unit to your configuration:
```bash
weather config --city "London" --unit fahrenheit
```
View your current settings:
```bash
weather config
```

### Options & Flags
- `-u`, `--unit`: Override the default temperature unit (`celsius`, `fahrenheit`, or `kelvin`) for a single command.
- `-c`, `--city`: Specify the city name.
- `-d`, `--days`: Specify the forecast length (1 to 5 days, defaults to 3).

---

## Data Providers
- **Weather Data & Geocoding**: [Open-Meteo API](https://open-meteo.com/) (Free for non-commercial use, no API keys needed).
- **IP Geolocation**: [ip-api.com](http://ip-api.com/) (Free geolocation service).

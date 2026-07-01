package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	unitFlag string
)

var rootCmd = &cobra.Command{
	Use:   "weather",
	Short: "Weather CLI tool to fetch current conditions and forecasts",
	Long:  `A fast and beautiful command-line weather client that fetches data from Open-Meteo.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing weather command: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&unitFlag, "unit", "u", "", "Temperature unit (celsius, fahrenheit, kelvin)")
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func BuildCli() *cobra.Command {
	var configPath string

	rootCmd := &cobra.Command{
		Use:   "hako",
		Short: "A key-value storage over HTTP.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Nothing to do. Use `hako --help` for help.")
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Start the server.",
		Run: func(cmd *cobra.Command, args []string) {
			PrepareLogger()
			var config Config

			if configPath != "" {
				result, err := LoadConfig(configPath)
				if err != nil {
					fmt.Printf("Failed to parse configuration file: %s\n", err.Error())
					os.Exit(1)
				}
				config = result
			} else {
				config = GetDefaultConfig()
			}

			StartServer(&config)
		},
	}

	runCmd.Flags().StringVar(&configPath, "config", "", "A path to configuration file")

	rootCmd.AddCommand(runCmd)

	return rootCmd
}

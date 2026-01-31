package main

import (
	"fmt"
	"os"

	"github.com/kostya-zero/hako/internal/config"
	"github.com/kostya-zero/hako/internal/server"
	"github.com/kostya-zero/hako/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
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
			utils.PrepareLogger()
			var cfg config.Config

			if configPath != "" {
				result, err := config.LoadConfig(configPath)
				if err != nil {
					fmt.Printf("Failed to parse configuration file: %s\n", err.Error())
					os.Exit(1)
				}
				cfg = result
			} else {
				cfg = config.GetDefaultConfig()
			}

			server.StartServer(&cfg)
		},
	}

	runCmd.Flags().StringVar(&configPath, "config", "", "A path to configuration file")

	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("An error occured in CLI interactions")
		os.Exit(1)
	}
}

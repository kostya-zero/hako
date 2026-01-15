package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func BuildCli() *cobra.Command {
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
			StartServer(&Config{Port: 3000, File: ""})
		},
	}

	rootCmd.AddCommand(runCmd)

	return rootCmd
}

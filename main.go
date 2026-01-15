package main

import (
	"fmt"
	"os"
)

func main() {
	cli := BuildCli()
	if err := cli.Execute(); err != nil {
		fmt.Printf("An error occured in CLI interactions")
		os.Exit(1)
	}
}

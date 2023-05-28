package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var NEBULOOM_LOGO = `
    _   __     __          __
   / | / /__  / /_  __  __/ /___  ____  ____ ___
  /  |/ / _ \/ __ \/ / / / / __ \/ __ \/ __ '__ \
 / /|  /  __/ /_/ / /_/ / / /_/ / /_/ / / / / / /
/_/ |_/\___/_.___/\__,_/_/\____/\____/_/ /_/ /_/
`

func main() {
	var rootCmd = &cobra.Command{
		Use:     "nebuloom",
		Short:   "Nebuloom",
		Version: "v0.1.0",
		Long:    "Nebuloom CLI tool to run Nebuloom files.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(NEBULOOM_LOGO)
			fmt.Println(`Nebuloom CLI ` + cmd.Version + "\n")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages in vendor directory using dep tool.",
	PreRun: func(cmd *cobra.Command, args []string) {
		check("dep")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var verboseFlag string
		if *verbose {
			verboseFlag = "-v"
		}
		o, e, err := run("dep", "ensure", verboseFlag)
		if err != nil {
			color.Red("\ninstall failed: %v\n\n", err)
			if *verbose {
				color.Red(e)
			}
			os.Exit(1)
		}
		if *verbose {
			color.White(o)
		}
		color.Green("\n✓ packages installed\n\n")
	},
}

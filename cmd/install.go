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
		check(DEP)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var verboseFlag string
		if *VERBOSE {
			verboseFlag = "-v"
		}
		o, e, err := run(DEP, "ensure", verboseFlag)
		if err != nil {
			color.Red("\ninstall failed: %v\n\n", err)
			if *VERBOSE {
				color.Red(e)
			}
			os.Exit(1)
		}
		if *VERBOSE {
			color.White(o)
		}
		color.Green("\n✓ packages installed\n\n")
	},
}

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
		checkFatal(DEP)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var v string
		if *VERBOSE {
			v = "-v"
		}

		_, err := run(DEP, "ensure", v)
		if err != nil {
			color.Red("\ninstall failed: %v\n\n", err)
			os.Exit(1)
		}
		color.Green("\n✓ packages installed\n\n")
	},
}

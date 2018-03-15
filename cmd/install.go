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
		a := []string{"ensure"}
		if *VERBOSE {
			a = append(a, "-v")
		}

		_, err := run(DEP, a...)
		if err != nil {
			color.Red("\ninstall failed: %v\n\n", err)
			os.Exit(1)
		}
		color.Green("\n✓ packages installed\n\n")
	},
}

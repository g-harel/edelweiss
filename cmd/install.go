package cmd

import (
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

		log("Installing/updating dependencies")
		_, err := run(DEP, a...)
		fatal(err)("install failed")

		color.Green("\n✓ packages installed\n\n")
	},
}

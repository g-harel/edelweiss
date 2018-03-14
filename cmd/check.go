package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func checkFatal(commands ...string) {
	for _, c := range commands {
		_, err := run(c)
		if err != nil {
			color.Red("\ndependency missing: %v: %v\n\n", c, err)
			os.Exit(1)
		}
	}
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Make sure all executable dependencies are in the path.",
	Run: func(cmd *cobra.Command, args []string) {
		checkFatal(GO, DEP, KUBECTL, MINIKUBE)
		color.Green("\n✓ all dependencies located\n\n")
	},
}

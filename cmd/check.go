package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func check(commands ...string) error {
	for _, c := range commands {
		_, _, err := run(c)
		if err != nil {
			return err
		}
	}
	return nil
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Make sure all executable dependencies are in the path.",
	Run: func(cmd *cobra.Command, args []string) {
		err := check(GO, DEP, KUBECTL, MINIKUBE)
		if err != nil {
			color.Red("\ndependency missing: %v\n\n", err)
			os.Exit(1)
		}
		color.Green("\n✓ all dependencies located\n\n")
	},
}

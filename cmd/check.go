package cmd

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func checkFatal(commands ...string) {
	for _, c := range commands {
		_, err := exec.LookPath(c)
		if err != nil {
			color.Red("\ncheckFatal: %s\n(If you are on windows, executables must have a file extension to be found)\n\n", err)
			os.Exit(1)
		}

		if *VERBOSE {
			color.White("found %s", c)
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

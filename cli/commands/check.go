package commands

import (
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func checkFatal(commands ...string) {
	for _, c := range commands {
		p, err := exec.LookPath(c)
		fatal(err)("check failed (If you are on windows, executables must have a file extension to be found)")
		verboseLog(color.White, "found %s (%s)", c, p)
	}
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Make sure all executable dependencies are in the path.",
	Run: func(cmd *cobra.Command, args []string) {
		log("Looking for all required executables")
		checkFatal(DOCKER, GO, HELM, KUBECTL)
		color.Green("\n✓ all dependencies located\n\n")
	},
}

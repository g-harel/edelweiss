package commands

import (
	"github.com/g-harel/edelweiss/cli"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Make sure all executable dependencies are in the path",
	Run: func(cmd *cobra.Command, args []string) {
		log.Progress("Looking for all required executables")
		cli.CheckFatal(DOCKER, GO, HELM, KUBECTL)

		log.Success("all dependencies located")
	},
}

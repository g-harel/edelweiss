package commands

import (
	"github.com/g-harel/edelweiss/cli"
	"github.com/spf13/cobra"
)

// TODO use clients for everything
// executable dependencies
const (
	GO      = "go"
	HELM    = "helm"
	KUBECTL = "kubectl"
)

var log = cli.Logger

var rootCmd = &cobra.Command{
	Use: "edelweiss",
}

func init() {
	cli.Flags(rootCmd)
}

// Execute executes the root command.
func Execute() {
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(installCmd)

	rootCmd.Execute()
}

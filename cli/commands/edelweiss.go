package commands

import (
	"github.com/g-harel/edelweiss/cli"
	"github.com/spf13/cobra"
)

var clilog = cli.Logger

var rootCmd = &cobra.Command{
	Use: "edelweiss",
}

func init() {
	cli.InitFlags(rootCmd)
}

// Execute executes the root command.
func Execute() {
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(installCmd)

	rootCmd.Execute()
}

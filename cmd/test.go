package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [services...]",
	Short: "Run go tests.",
	Long:  "Run go tests. Service names can be passed as arguments. Providing no names will run all tests.",
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue(*WORKDIR)
		color.Blue(strings.Join(args, "\n"))
		o, e, err := run(GO, "test", "-race", "-v", "./..")
		if err != nil {
			color.Red("\ntests failed: %v\n\n", err)
			if *VERBOSE {
				color.Red(e)
			}
			os.Exit(1)
		}
		if *VERBOSE {
			color.White(o)
		}
		color.Green("\n✓ checks passed\n\n")
	},
}

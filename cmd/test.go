package cmd

import (
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [services...]",
	Short: "Run go tests.",
	Long:  "Run go tests. Service names can be passed as arguments. Providing no names will run all tests.",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkFatal(GO)
	},
	Run: func(cmd *cobra.Command, args []string) {
		dirs := []string{"./..."}
		if len(args) > 0 {
			dirs = make([]string, len(args))
			for i, s := range args {
				dirs[i] = path.Join("./services", s)
			}
		}

		a := []string{"test", "./...", "-race"}
		if *VERBOSE {
			a = append(a, "-v")
		}

		for _, d := range dirs {
			out, err := run(GO, append(a, d)...)
			if err != nil {
				log(color.White, out)
				color.Red("\ntest failed: %v\n\n", err)
				os.Exit(1)
			}
		}
		color.Green("\n✓ tests passed\n\n")
	},
}
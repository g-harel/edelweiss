package commands

import (
	"path"

	"github.com/g-harel/edelweiss/cli"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [services...]",
	Short: "Run go tests",
	Long:  "Run go tests. Service names can be passed as arguments. Providing no names will run all tests.",
	Run: func(cmd *cobra.Command, args []string) {
		dirs := []string{"./..."}
		if len(args) > 0 {
			dirs = make([]string, len(args))
			for i, s := range args {
				dirs[i] = path.Join("./services", s)
			}
		}

		a := []string{"test", "./...", "-race"}
		if *cli.VERBOSE {
			a = append(a, "-v")
		}

		for _, d := range dirs {
			clilog.Progress("Running tests in \"%v\" dir", d)
			out, err := cli.Exec("go", append(a, d)...)
			clilog.Fatal(err, "test failed: %v", out)
		}

		clilog.Success("tests passed")
	},
}

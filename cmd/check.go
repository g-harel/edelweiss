package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func checkFatal(commands ...string) {
	path, ok := os.LookupEnv("PATH")
	if !ok {
		color.Red("\ncould not read PATH\n\n")
		os.Exit(1)
	}

	dirs := strings.Split(path, string(os.PathListSeparator))

	exec := make(map[string]bool)
	for _, d := range dirs {
		files, err := ioutil.ReadDir(d)
		if err != nil {
			color.Red("\ncould not read directory (%v): %v\n\n", d, err)
			os.Exit(1)
		}

		for _, f := range files {
			if f.Mode()&0111 != 0 {
				exec[f.Name()] = true
			}
		}
	}

	for _, c := range commands {
		if !exec[c] {
			color.Red("\ndependency missing: %v\n\n", c)
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

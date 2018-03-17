package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// executable dependencies
const (
	DEP      = "dep"
	GO       = "go"
	HELM     = "helm"
	KUBECTL  = "kubectl"
	MINIKUBE = "minikube"
)

// global flags
var (
	VERBOSE *bool
	WORKDIR *string
)

func run(command string, args ...string) (string, error) {
	log(color.Yellow, "running: %v %v", command, strings.Join(args, " "))
	b, err := exec.Command(command, args...).CombinedOutput()
	log(color.White, "%s\n", b)
	return string(b), err
}

func fatal(err error) func(format string, args ...interface{}) {
	return func(f string, a ...interface{}) {
		if err != nil {
			color.Red("\n%s: %v\n\n", fmt.Sprintf(f, a...), err)
			os.Exit(1)
		}
	}
}

func log(f func(string, ...interface{}), format string, args ...interface{}) {
	if *VERBOSE {
		f(format, args...)
	}
}

var rootCmd = &cobra.Command{
	Use: "edelweiss",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiWhite("Hello World!")
	},
}

func init() {
	dir, err := os.Getwd()
	fatal(err)("Could not locate working directory")

	VERBOSE = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output and logging")
	WORKDIR = rootCmd.PersistentFlags().StringP("workdir", "w", dir, "define working directory")
}

// Execute executes the root command.
func Execute() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(rookCmd)

	rootCmd.Execute()
}
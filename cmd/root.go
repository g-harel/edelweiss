package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// executable dependencies
const (
	GO       = "go"
	DEP      = "dep"
	MINIKUBE = "minikube"
	KUBECTL  = "kubectl"
)

// global flags
var (
	VERBOSE *bool
	WORKDIR *string
)

func run(command string, args ...string) (string, error) {
	if *VERBOSE {
		color.Yellow("running: %v %v", command, strings.Join(args, " "))
	}

	b, err := exec.Command(command, args...).CombinedOutput()

	if *VERBOSE {
		color.White("%s\n", b)
	}

	return string(b), err
}

var rootCmd = &cobra.Command{
	Use: "edelweiss",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiWhite("Hello World!")
	},
}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		color.Red("\nCould not locate working directory: %v\n\n", err)
	}

	VERBOSE = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output and logging")
	WORKDIR = rootCmd.PersistentFlags().StringP("workdir", "w", dir, "define working directory")
}

// Execute executes the root command.
func Execute() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(testCmd)

	rootCmd.Execute()
}

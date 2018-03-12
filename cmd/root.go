package cmd

import (
	"bytes"
	"fmt"
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

func run(command string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(command, args...)

	if *VERBOSE {
		color.Yellow("running: %v %v", command, strings.Join(args, " "))
	}

	outp, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	errp, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", "", err
	}

	bo := new(bytes.Buffer)
	bo.ReadFrom(outp)

	be := new(bytes.Buffer)
	be.ReadFrom(errp)

	err = cmd.Wait()
	if err != nil {
		return "", "", nil
	}

	return bo.String(), bo.String(), nil
}

var rootCmd = &cobra.Command{
	Use: "edelweiss",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
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

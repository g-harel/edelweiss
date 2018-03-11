package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func run(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	b.ReadFrom(stdout)

	err = cmd.Wait()
	if err != nil {
		return "", nil
	}

	return b.String(), nil
}

var rootCmd = &cobra.Command{
	Use: "edelweiss",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

var checkCmd = &cobra.Command{
	Use: "check",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"go", "dep", "kubectl", "minikube"}
		for _, c := range commands {
			_, err := run(c)
			if err != nil {
				color.Red("\ndependency missing: %v\n\n", err)
				os.Exit(1)
			}
		}
		color.Green("\n✓ all dependencies located\n\n")
	},
}

// Execute executes the root command.
func Execute() {
	rootCmd.AddCommand(checkCmd)

	rootCmd.Execute()
}

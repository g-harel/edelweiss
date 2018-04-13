package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// global cli flags
var (
	VERBOSE *bool
	WORKDIR *string
)

// InitFlags initializes the values of the global flags.
func InitFlags(cmd *cobra.Command) {
	dir, err := os.Getwd()
	Logger.Fatal(err, "Could not locate working directory")

	VERBOSE = cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output and logging")
	WORKDIR = cmd.PersistentFlags().StringP("workdir", "w", dir, "define working directory")
}

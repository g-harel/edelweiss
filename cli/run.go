package cli

import (
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// Run runs the provided executable (if installed).
func Run(command string, args ...string) (string, error) {
	Logger.DebugC(color.Yellow, "$ %v %v", command, strings.Join(args, " "))
	b, err := exec.Command(command, args...).CombinedOutput()
	Logger.Debug("%s\n", b)
	return string(b), err
}

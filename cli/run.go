package cli

import (
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Run(command string, args ...string) (string, error) {
	Logger.DebugC(color.Yellow, "$ %v %v", command, strings.Join(args, " "))
	b, err := exec.Command(command, args...).CombinedOutput()
	Logger.Debug("%s\n", b)
	return string(b), err
}

package cli

import "os/exec"

func CheckFatal(commands ...string) {
	for _, c := range commands {
		p, err := exec.LookPath(c)
		Logger.Fatal(err, "check failed (If you are on windows, executables must have a file extension to be found)")
		Logger.Debug("found %s (%s)", c, p)
	}
}

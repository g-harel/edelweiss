package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Logger is a shared instance of the cli's logger.
// This allows other cli packages to consume the same global logger.
var Logger = logger{}

type logger struct{}

func (l *logger) Progress(format string, args ...interface{}) {
	color.New(color.Bold, color.FgHiBlue).Printf("> ")
	color.White("%s...\n", fmt.Sprintf(format, args...))
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.DebugC(color.White, format, args...)
}

func (l *logger) DebugC(f func(string, ...interface{}), format string, args ...interface{}) {
	if *VERBOSE {
		f(format, args...)
	}
}

func (l *logger) Success(format string, args ...interface{}) {
	color.Green("\n✓ %s\n\n", fmt.Sprintf(format, args...))
}

func (l *logger) Fatal(err error, format string, args ...interface{}) {
	if err != nil {
		color.Red("\n%s: %v\n\n", fmt.Sprintf(format, args...), err)
		os.Exit(1)
	}
}

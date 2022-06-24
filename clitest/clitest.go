package clitest

import (
	"bytes"
	"time"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/internal/types"
)

// Bool returns a bool to bind flags/args to.
func Bool(v bool) bool {
	return types.NewBool(nil, v).Get()
}

// String returns a string to bind flags/args to.
func String(v string) string {
	return types.NewString(nil, v).Get()
}

// Int returns a int to bind flags/args to.
func Int(v int) int {
	return types.NewInt(nil, v).Get()
}

// Float64 returns a float64 to bind flags/args to.
func Float64(v float64) float64 {
	return types.NewFloat64(nil, v).Get()
}

// Duration returns a time.Duration to bind flags/args to.
func Duration(v time.Duration) time.Duration {
	return types.NewDuration(nil, v).Get()
}

// StringSlice returns a []string to bind flags/args to.
func StringSlice(v []string) []string {
	return types.NewStringSlice(nil, v).Get()
}

// Run runs the test command.
func Run(cmd *Command, args ...string) (string, error) {
	var buf bytes.Buffer

	cmd.SetOutput(&buf)

	if err := cli.Run(makeArgs(cmd, args...), cmd); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// makeArgs returns a slice of arguments constructed from the full name of the
// command and the provided arguments.
func makeArgs(root *Command, args ...string) []string {
	allArgs := make([]string, 0)
	allArgs = append(allArgs, root.cmd.FullName())

	for _, arg := range args {
		allArgs = append(allArgs, arg)
	}

	return allArgs
}

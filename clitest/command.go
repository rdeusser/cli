package clitest

import (
	"io"

	"github.com/rdeusser/cli"
)

// Command is a wrapper around a command and a run function to simulate a
// runner.
type Command struct {
	cmd   cli.Command
	runFn RunFn
}

// RunFn represents the `Run() error` method of a command to be used in tests.
type RunFn func() error

// NewCommand creates a new command.
func NewCommand(cmd cli.Command, runFn RunFn) *Command {
	return &Command{
		cmd:   cmd,
		runFn: runFn,
	}
}

// Init returns the wrapped command.
func (c *Command) Init() cli.Command {
	return c.cmd
}

// Run runs the commands `Run() error` method.
//
// If the commands runFn is nil, it returns an error indicating that the help
// output should be returned.
func (c *Command) Run() error {
	if c.runFn == nil {
		return cli.ErrPrintHelp
	}

	return c.runFn()
}

// AddCommands adds commands to the test command.
func (c *Command) AddCommands(commands ...*Command) {
	for _, command := range commands {
		c.cmd.AddCommands(
			NewCommand(command.cmd, command.runFn),
		)
	}
}

// SetOutput sets the output used by the commmand.
func (c *Command) SetOutput(w io.Writer) {
	c.cmd.SetOutput(w)
}

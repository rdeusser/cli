package cli_test

import (
	"fmt"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/clitest"
)

func TestCommand(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			"Root command output",
			nil,
		},
		{
			"Root command help",
			[]string{"--help"},
		},
		{
			"Server start subcommand output",
			[]string{"server", "start"},
		},
		{
			"Server start subcommand help",
			[]string{"server", "start", "--help"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootCommand := clitest.NewCommand(cli.Command{
				Name: "test",
				Desc: "A test binary",
			}, nil)

			serverCommand := clitest.NewCommand(cli.Command{
				Name: "server",
				Desc: "Do something with a test server",
			}, nil)

			startCommand := clitest.NewCommand(cli.Command{
				Name: "start",
				Desc: "Starts a test server",
			}, func() error {
				fmt.Fprintln(rootCommand.Output(), "running from server start subcommand")
				return nil
			})

			serverCommand.AddCommands(
				startCommand,
			)

			rootCommand.AddCommands(
				serverCommand,
			)

			output, err := clitest.Run(rootCommand, tc.args...)
			assert.NoError(t, err)
			assert.NotEmpty(t, output)

			autogold.Equal(t, autogold.Raw(output))
		})
	}
}

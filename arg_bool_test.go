package cli_test

import (
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/clitest"
)

func TestBoolArg(t *testing.T) {
	testCases := []struct {
		name     string
		bind     bool
		args     []string
		expected bool
	}{
		{
			"Help output",
			clitest.Bool(false),
			[]string{"--help"},
			false,
		},
		{
			"Set value (false)",
			clitest.Bool(false),
			[]string{"false"},
			false,
		},
		{
			"Set value (true)",
			clitest.Bool(false),
			[]string{"true"},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			arg := &cli.BoolArg{
				Bind:     &tc.bind,
				Name:     "test-bool-arg",
				Desc:     "test bool",
				Position: 0,
			}

			cmd := clitest.NewCommand(cli.Command{
				Name: "test-bool-arg",
				Desc: "Test bool arg",
				Args: cli.Args{
					arg,
				},
			}, nil)

			output, err := clitest.Run(cmd, tc.args...)
			assert.NoError(t, err)
			assert.NotEmpty(t, output)

			autogold.Equal(t, autogold.Raw(output))
		})
	}
}

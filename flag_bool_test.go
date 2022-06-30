package cli_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/clitest"
)

func TestBoolFlag(t *testing.T) {
	testCases := []struct {
		name     string
		bind     bool
		envVar   string
		args     []string
		expected bool
	}{
		{
			"Help output",
			clitest.Bool(false),
			"",
			[]string{"--help"},
			false,
		},
		{
			"Default value (false)",
			clitest.Bool(false),
			"",
			nil,
			false,
		},
		{
			"Default value (true)",
			clitest.Bool(true),
			"",
			nil,
			true,
		},
		{
			"Default value (envvar)",
			clitest.Bool(false),
			"TEST",
			nil,
			true,
		},
		{
			"Set value using shorthand",
			clitest.Bool(false),
			"",
			[]string{"-t"},
			true,
		},
		{
			"Set value using name",
			clitest.Bool(false),
			"",
			[]string{"--test"},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag := &cli.BoolFlag{
				Bind:      &tc.bind,
				Name:      "test",
				Shorthand: "t",
				Desc:      "test flag",
				Default:   tc.bind,
			}

			if tc.envVar != "" {
				flag.EnvVar = tc.envVar
				os.Setenv(flag.EnvVar, "true")
			}

			cmd := clitest.NewCommand(cli.Command{
				Name: "test-bool-flag",
				Desc: "Test bool flag",
				Flags: cli.Flags{
					flag,
				},
			}, func() error {
				fmt.Fprintln(clitest.Output, flag.String())
				return nil
			})

			output, err := clitest.Run(cmd, tc.args...)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, *flag.Bind)

			autogold.Equal(t, autogold.Raw(output))
		})
	}
}

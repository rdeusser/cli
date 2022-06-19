package cli

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli/internal/types"
)

type testCommand struct {
	Debug bool
}

func (cmd *testCommand) Init() Command {
	return Command{
		Name: "test",
		Desc: "test setting a bool arg",
		Args: Args{
			&BoolArg{
				Bind:     &cmd.Debug,
				Name:     "debug",
				Desc:     "debug mode",
				Position: 0,
				Required: true,
			},
		},
	}
}

func (testCommand) Run() error {
	return nil
}

func TestBoolArg(t *testing.T) {
	NoColor = true // autogold seems to have problems with color in golden files

	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			"Set value",
			[]string{"true"},
			true,
		},
		{
			"Set value with non-bool",
			[]string{"foo"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err    error
				runner boolCommand
			)

			cmd := runner.Init()

			var buf bytes.Buffer
			Output = &buf

			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()

			os.Args = []string{}
			os.Args = append(os.Args, cmd.Name)
			os.Args = append(os.Args, tt.args...)

			cmd, err = Run(&testCommand{})
			assert.NoError(t, err)

			assert.Equal(t, types.BoolType, cmd.args[0].Type())
			assert.Equal(t, tt.expected, cmd.args[0].Bind)
			assert.Equal(t, strconv.FormatBool(tt.expected), cmd.args[0].Bind)

			autogold.Equal(t, autogold.Raw(buf.String()))
		})
	}
}

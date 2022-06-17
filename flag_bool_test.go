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

var boolFlag = &BoolFlag{
	Name:      "test",
	Shorthand: "t",
	Desc:      "run tests",
	Default:   false,
	Required:  true,
}

type boolCommand struct{}

func (boolCommand) Init() Command {
	return Command{
		Name: "test",
		Desc: "test setting a bool flag",
		Flags: []Flag{
			boolFlag,
		},
	}
}

func (boolCommand) Run() error {
	return nil
}

func TestBoolFlag(t *testing.T) {
	NoColor = true // autogold seems to have problems with color in golden files

	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			"Help Output",
			[]string{"--help"},
			false,
		},
		{
			"Default Value",
			[]string{""},
			false,
		},
		{
			"Set Value Using Shorthand",
			[]string{"-t"},
			true,
		},
		{
			"Set Value Using Name",
			[]string{"--test"},
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

			cmd, err = Run(&boolCommand{})
			assert.NoError(t, err)

			assert.Equal(t, types.BoolType, boolFlag.Type())
			assert.Equal(t, tt.expected, boolFlag.Get())
			assert.Equal(t, strconv.FormatBool(tt.expected), boolFlag.String())

			autogold.Equal(t, autogold.Raw(buf.String()))
		})
	}
}

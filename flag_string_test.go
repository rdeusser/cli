package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stringFlag = &StringFlag{
	Name:      "filename",
	Shorthand: "f",
	Desc:      "specify filename",
	Default:   "config.yaml",
	Required:  true,
}

type stringCommand struct{}

func (stringCommand) Init() Command {
	return Command{
		Name: "test",
		Desc: "test setting a string flag",
		Flags: []Flag{
			stringFlag,
		},
	}
}

func (stringCommand) Run(args []string) error {
	return nil
}

func TestStringFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		snapshot bool
	}{
		{
			"Help Output",
			[]string{"--help"},
			"config.yaml",
			true,
		},
		{
			"Default Value",
			[]string{""},
			"config.yaml",
			false,
		},
		{
			"Set Value Using Shorthand",
			[]string{"-f", "anotherconfig.yaml"},
			"anotherconfig.yaml",
			false,
		},
		{
			"Set Value Using Name",
			[]string{"--filename", "anotherconfig.yaml"},
			"anotherconfig.yaml",
			false,
		},
		{
			"Set Value After Args",
			[]string{"hi", "-f", "config.yaml"},
			"config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err    error
				runner stringCommand
			)

			cmd := runner.Init()

			var buf bytes.Buffer
			SetOutput(&buf)

			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()

			os.Args = []string{}
			os.Args = append(os.Args, cmd.Name)
			os.Args = append(os.Args, tt.args...)

			cmd, err = Run(&stringCommand{})
			assert.NoError(t, err)

			assert.Equal(t, String, stringFlag.Type())
			assert.Equal(t, tt.expected, stringFlag.Get())
			assert.Equal(t, tt.expected, stringFlag.String())

			if tt.snapshot {
				snapshot(t, buf.Bytes(), nil)
			}
		})
	}
}

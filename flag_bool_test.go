package cli

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

func (boolCommand) Run(args []string) error {
	return nil
}

func TestBoolFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
		snapshot bool
	}{
		{
			"Help Output",
			[]string{"--help"},
			false,
			true,
		},
		{
			"Default Value",
			[]string{""},
			false,
			false,
		},
		{
			"Set Value Using Shorthand",
			[]string{"-t"},
			true,
			false,
		},
		{
			"Set Value Using Name",
			[]string{"--test"},
			true,
			false,
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
			SetOutput(&buf)

			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()

			os.Args = []string{}
			os.Args = append(os.Args, cmd.Name)
			os.Args = append(os.Args, tt.args...)

			cmd, err = Run(&boolCommand{})
			assert.NoError(t, err)

			assert.Equal(t, Bool, boolFlag.Type())
			assert.Equal(t, tt.expected, boolFlag.Get())
			assert.Equal(t, strconv.FormatBool(tt.expected), boolFlag.String())

			if tt.snapshot {
				snapshot(t, buf.Bytes(), nil)
			}
		})
	}
}

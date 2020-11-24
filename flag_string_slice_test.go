package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stringSliceFlag = &StringSliceFlag{
	Name:      "tests",
	Shorthand: "t",
	Desc:      "run tests",
	Default:   []string{"pkg1", "pkg2"},
	Required:  true,
}

type stringSliceCommand struct{}

func (stringSliceCommand) Init() Command {
	return Command{
		Name: "test",
		Desc: "test setting a bool flag",
		Flags: []Flag{
			stringSliceFlag,
		},
	}
}

func (stringSliceCommand) Run(args []string) error {
	return nil
}

func TestStringSliceFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
		snapshot bool
	}{
		{
			"Help Output",
			[]string{"--help"},
			[]string{"pkg1", "pkg2"},
			true,
		},
		{
			"Default Value",
			[]string{""},
			[]string{"pkg1", "pkg2"},
			false,
		},
		{
			"Set Value Using Shorthand",
			[]string{"-t", "pkg3", "pkg4"},
			[]string{"pkg3", "pkg4"},
			false,
		},
		{
			"Set Value Using Name",
			[]string{"--tests", "pkg3", "pkg4"},
			[]string{"pkg3", "pkg4"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err    error
				runner stringSliceCommand
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

			cmd, err = Run(&stringSliceCommand{})
			assert.NoError(t, err)

			assert.Equal(t, StringSlice, stringSliceFlag.Type())
			assert.Equal(t, tt.expected, stringSliceFlag.Get())

			if tt.snapshot {
				snapshot(t, buf.Bytes(), nil)
			}
		})
	}
}

package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

const testdata = "testdata"

type rootCommand struct{}

func (rootCommand) Init() Command {
	cmd := Command{
		Name: "test",
		Desc: "Test CLI",
	}

	cmd.AddCommands(
		&serverCommand{},
	)

	return cmd
}

func (rootCommand) Run(args []string) error {
	fmt.Fprintln(output, "running from the root command")
	return nil
}

type serverCommand struct{}

func (serverCommand) Init() Command {
	cmd := Command{
		Name: "server",
		Desc: "Runs a server for the test command (not really)",
	}

	cmd.AddCommands(
		&serverStartCommand{},
	)

	return cmd
}

func (serverCommand) Run(args []string) error {
	fmt.Fprintln(output, "running from server subcommand")
	return nil
}

type serverStartCommand struct{}

func (serverStartCommand) Init() Command {
	return Command{
		Name: "start",
		Desc: "Starts a server of some sort",
	}
}

func (serverStartCommand) Run(args []string) error {
	fmt.Fprintln(output, "running from server start subcommand")
	return nil
}

func snapshot(t *testing.T, output []byte, obj interface{}) {
	t.Helper()

	testname := strings.ReplaceAll(t.Name(), "/", "-")
	filename := filepath.Join(testdata, testname)
	snapshot, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			os.Setenv("UPDATE_SNAPSHOTS", "true")
			t.Logf("Snapshot %s doesn't exist yet", testname)
		} else {
			assert.NoError(t, err)
		}
	}

	if _, ok := os.LookupEnv("UPDATE_SNAPSHOTS"); ok {
		if obj != nil {
			err := ioutil.WriteFile(filename, []byte(pretty.Sprint(obj)), os.FileMode(0644))
			assert.NoError(t, err)
		} else {
			err := ioutil.WriteFile(filename, output, os.FileMode(0644))
			assert.NoError(t, err)
		}

		t.Logf("Updated snapshot for %s", testname)
	} else {
		if obj != nil {
			assert.Equal(t, snapshot, pretty.Sprint(obj))
		} else {
			assert.Equal(t, snapshot, output)
		}
	}
}

func TestCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			"Root Command Output",
			[]string{""},
		},
		{
			"Root Command Help",
			[]string{"--help"},
		},
		{
			"Server Subcommand Output",
			[]string{"server"},
		},
		{
			"Server Subcommand Help",
			[]string{"server", "--help"},
		},
		{
			"Server Start Subcommand Output",
			[]string{"server", "start"},
		},
		{
			"Server Start Subcommand Help",
			[]string{"server", "start", "--help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err    error
				runner rootCommand
			)

			cmd := runner.Init()

			var buf bytes.Buffer
			SetOutput(&buf)

			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()

			os.Args = []string{""}
			os.Args = append(os.Args, cmd.Name)
			os.Args = append(os.Args, tt.args...)

			cmd, err = Run(&rootCommand{})
			assert.NoError(t, err)

			snapshot(t, buf.Bytes(), nil)
		})
	}
}

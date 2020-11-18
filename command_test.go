package cli

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

var testdata = cupaloy.New(
	cupaloy.SnapshotSubdirectory("testdata"),
	cupaloy.ShouldUpdate(func() bool {
		return true
	}),
)

type rootCommand struct{}

func (rootCommand) Init() Command {
	cmd := Command{
		Name: "test",
		Desc: "Test CLI",
	}

	cmd.AddCommands(
		&serverCommand{},
		&cacheCacheCommand{},
	)

	return cmd
}

func (rootCommand) Run() error {
	fmt.Println("running from the root command")
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

func (serverCommand) Run() error {
	fmt.Println("running from server subcommand")
	return nil
}

type serverStartCommand struct{}

func (serverStartCommand) Init() Command {
	return Command{
		Name: "start",
		Desc: "Starts a server of some sort",
	}
}

func (serverStartCommand) Run() error {
	fmt.Println("running from server start subcommand")
	return nil
}

type cacheCacheCommand struct{}

func (cacheCacheCommand) Init() Command {
	return Command{
		Name: "cachecache",
		Desc: "Does something cache like",
	}
}

func (cacheCacheCommand) Run() error {
	fmt.Println("running from cachecache subcommand")
	return nil
}

func TestCommand_Help(t *testing.T) {
	var root rootCommand

	cmd := root.Init()

	var buf bytes.Buffer
	SetOutput(&buf)

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{cmd.Name, "--help"}

	err := Run(&rootCommand{})
	assert.NoError(t, err)

	testdata.SnapshotT(t, buf.String())
}

func TestSubCommand_Help(t *testing.T) {
	var root rootCommand

	cmd := root.Init()

	var buf bytes.Buffer
	SetOutput(&buf)

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{cmd.Name, "server", "--help"}

	err := Run(&rootCommand{})
	assert.NoError(t, err)

	testdata.SnapshotT(t, buf.String())
}

func TestSubSubCommand_Help(t *testing.T) {
	var root rootCommand

	cmd := root.Init()

	buf := &bytes.Buffer{}
	SetOutput(buf)

	oldArgs := os.Args

	os.Args = []string{cmd.Name, "server", "start", "--help"}
	defer func() {
		os.Args = oldArgs
	}()

	err := Run(&rootCommand{})
	assert.NoError(t, err)

	testdata.SnapshotT(t, buf.String())
}

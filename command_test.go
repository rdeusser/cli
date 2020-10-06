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
)

type serverCommand struct{}

func (c *serverCommand) Init() Command {
	return Command{
		Name: "server",
		Desc: "Runs a server for the test command (not really)",
	}
}

func (c *serverCommand) Run() error {
	fmt.Println("running from server subcommand")
	return nil
}

type serverStartCommand struct{}

func (c *serverStartCommand) Init() Command {
	return Command{
		Name: "start",
		Desc: "Starts a server of some sort",
	}
}

func (c *serverStartCommand) Run() error {
	fmt.Println("running from server start subcommand")
	return nil
}

type cacheCacheCommand struct{}

func (c *cacheCacheCommand) Init() Command {
	return Command{
		Name: "cachecache",
		Desc: "Does something cache like",
	}
}

func (c *cacheCacheCommand) Run() error {
	fmt.Println("running from cachecache subcommand")
	return nil
}

func TestCommand_Help(t *testing.T) {
	cmd := &Command{
		Name: "test",
		Desc: "Test CLI",
	}

	cmd.AddCommands(
		&serverCommand{},
		&cacheCacheCommand{},
	)

	var buf bytes.Buffer
	cmd.SetOutput(&buf)

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{cmd.Name, "--help"}

	err := cmd.Run()
	assert.NoError(t, err)

	err = testdata.Snapshot(buf.String())
	assert.NoError(t, err)
}

func TestSubCommand_Help(t *testing.T) {
	cmd := &Command{
		Name: "test",
		Desc: "Test CLI",
	}

	cmd.AddCommands(
		&serverCommand{},
		&cacheCacheCommand{},
	)

	var buf bytes.Buffer
	cmd.SetOutput(&buf)

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{cmd.Name, "server", "--help"}

	err := cmd.Run()
	assert.NoError(t, err)

	err = testdata.Snapshot(buf.String())
	assert.NoError(t, err)
}

func TestSubSubCommand_Help(t *testing.T) {
	var server serverCommand

	cmd := server.Init()

	cmd.AddCommands(
		&serverStartCommand{},
	)

	var buf bytes.Buffer
	cmd.SetOutput(&buf)

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{cmd.Name, "server", "start", "--help"}

	err := cmd.Run()
	assert.NoError(t, err)

	err = testdata.Snapshot(buf.String())
	assert.NoError(t, err)
}

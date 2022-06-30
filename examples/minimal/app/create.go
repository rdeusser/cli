package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type CreateCommand struct {
	Test   bool
	Server string
}

func (cmd *CreateCommand) Init() cli.Command {
	return cli.Command{
		Name: "create",
		Desc: "Creates some things, probably",
		Flags: cli.Flags{
			&cli.BoolFlag{
				Bind:      &cmd.Test,
				Name:      "test",
				Shorthand: "t",
				Desc:      "Create a test server",
				Default:   false,
				EnvVar:    "TEST",
				Required:  false,
			},
		},
		Args: cli.Args{
			&cli.StringArg{
				Bind:     &cmd.Server,
				Name:     "server",
				Desc:     "Server name",
				Required: true,
			},
		},
	}
}

func (cmd *CreateCommand) Run() error {
	fmt.Println("running from the create command")
	fmt.Println("test:   ", cmd.Test)
	fmt.Println("server: ", cmd.Server)
	return nil
}

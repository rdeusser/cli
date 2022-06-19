package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type CreateCommand struct {
	Server bool
}

func (cmd *CreateCommand) Init() cli.Command {
	return cli.Command{
		Name: "create",
		Desc: "Creates some things, probably",
		Args: cli.Args{
			&cli.BoolArg{
				Bind:     &cmd.Server,
				Name:     "server",
				Desc:     "Create server?",
				Position: 0,
				Required: true,
			},
		},
	}
}

func (cmd *CreateCommand) Run() error {
	fmt.Println("running from the create command")
	fmt.Println(cmd.Server)
	return nil
}

package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type CreateCommand struct{}

func (cmd *CreateCommand) Init() cli.Command {
	return cli.Command{
		Name: "create",
		Desc: "Creates some things, probably",
	}
}

func (cmd *CreateCommand) Run() error {
	fmt.Println("running from the create command")
	return nil
}

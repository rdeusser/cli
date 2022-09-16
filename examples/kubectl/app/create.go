package app

import "github.com/rdeusser/cli"

type CreateCommand struct{}

func (CreateCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "create",
		Desc: "Create a resource",
	}

	return cmd
}

func (CreateCommand) Run() error {
	return cli.ErrPrintHelp
}

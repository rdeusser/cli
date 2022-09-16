package app

import "github.com/rdeusser/cli"

type DeleteCommand struct{}

func (DeleteCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "delete",
		Desc: "Delete a resource",
	}

	return cmd
}

func (DeleteCommand) Run() error {
	return cli.ErrPrintHelp
}

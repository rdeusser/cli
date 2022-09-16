package app

import "github.com/rdeusser/cli"

type GetCommand struct{}

func (GetCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "get",
		Desc: "Get a resource",
	}

	cmd.AddCommands(
		&GetPodCommand{},
	)

	return cmd
}

func (GetCommand) Run() error {
	return cli.ErrPrintHelp
}

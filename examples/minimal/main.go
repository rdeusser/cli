package main

import (
	"log"
	"os"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/examples/minimal/app"
)

type Minimal struct{}

func (Minimal) Init() cli.Command {
	cmd := cli.Command{
		Name: "minimal",
		Desc: "Minimal example",
	}

	cmd.AddCommands(
		&app.CreateCommand{},
	)

	return cmd
}

func (Minimal) Run() error {
	return cli.ErrPrintHelp
}

func main() {
	if err := cli.Run(os.Args, &Minimal{}); err != nil {
		log.Fatal(err)
	}
}

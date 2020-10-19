package main

import (
	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/examples/minimal/app"
)

func main() {
	cmd := cli.Command{
		Name: "minimal",
		Desc: "Minimal example",
	}

	cmd.AddCommands(
		&app.CreateCommand{},
	)

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/examples/minimal/app"
)

type Minimal struct{}

func (Minimal) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "minimal",
		Desc: "Minimal example",
	}

	cmd.AddCommands(
		&app.CreateCommand{},
	)

	return cmd
}

func (Minimal) PersistentPreRun() error {
	fmt.Println("[main] running main persistent pre-runner")
	return nil
}

func (Minimal) PreRun() error {
	fmt.Println("[main] running main pre-runner")
	return nil
}

func (Minimal) Run() error {
	fmt.Println("[main] running main runner")
	return nil
}

func (Minimal) PostRun() error {
	fmt.Println("[main] running main post-runner")
	return nil
}

func (Minimal) PersistentPostRun() error {
	fmt.Println("[main] running main persistent post-runner")
	return nil
}

func main() {
	if err := cli.Execute(&Minimal{}, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

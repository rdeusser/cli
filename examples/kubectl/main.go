package main

import (
	"fmt"
	"os"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/cli/examples/kubectl/app"
)

type RootCommand struct {
	Debug     bool
	Namespace string
}

func (*RootCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "kubectl",
		Desc: "kubectl controls the Kubernetes cluster manager",
		Flags: cli.Flags{
			&cli.Flag[bool]{
				Name: "debug",
				Desc: "Set logging level to debug",
			},
			&cli.Flag[string]{
				Name:      "namespace",
				Shorthand: "n",
				Desc:      "Namespace to operate on",
				Default:   "default",
			},
			&cli.Flag[bool]{
				Shorthand: "A",
				Desc:      "All namespaces",
			},
		},
	}

	cmd.AddCommands(
		&app.ApplyCommand{},
		&app.CreateCommand{},
		&app.DeleteCommand{},
		&app.GetCommand{},
	)

	return cmd
}

func (*RootCommand) Run() error {
	return cli.ErrPrintHelp
}

func main() {
	if err := cli.Execute(&RootCommand{}, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

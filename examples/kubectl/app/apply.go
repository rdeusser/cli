package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type ApplyCommand struct {
	Filename cli.Path
}

func (ac *ApplyCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "apply",
		Desc: "Apply a resource",
		Flags: cli.Flags{
			&cli.Flag[cli.Path]{
				Name:      "filename",
				Shorthand: "f",
				Desc:      "Filename to apply",
				Value:     &ac.Filename,
				Required:  true,
			},
		},
	}

	return cmd
}

func (ac *ApplyCommand) Run() error {
	fmt.Printf("Filename: %s\n", &ac.Filename)
	return nil
}

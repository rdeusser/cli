package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type ApplyCommand struct {
	Debug    bool
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

// SetOptions is how we get flag values from parent commands. The root command
// has a debug flag. We don't want to copy and paste that flag to every command
// nor do we want to take that flag and put it in a different package so all the
// other commands can import it. We simply need the value. SetOptions lets us
// get that value.
func (ac *ApplyCommand) SetOptions(flags cli.Flags) error {
	ac.Debug = cli.ValueOf[bool](flags, "debug")

	return nil
}

func (ac *ApplyCommand) Run() error {
	fmt.Printf("Filename: %s\n", &ac.Filename)
	return nil
}

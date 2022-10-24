package app

import "github.com/rdeusser/cli"

type DeleteCommand struct {
	Debug bool
}

func (DeleteCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "delete",
		Desc: "Delete a resource",
	}

	return cmd
}

// SetOptions is how we get flag values from parent commands. The root command
// has a debug flag. We don't want to copy and paste that flag to every command
// nor do we want to take that flag and put it in a different package so all the
// other commands can import it. We simply need the value. SetOptions lets us
// get that value.
func (dc *DeleteCommand) SetOptions(flags cli.Flags) error {
	dc.Debug = cli.ValueOf[bool](flags, "debug")

	return nil
}

func (DeleteCommand) Run() error {
	return cli.ErrPrintHelp
}

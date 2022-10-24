package app

import "github.com/rdeusser/cli"

type CreateCommand struct {
	Debug bool
}

func (CreateCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "create",
		Desc: "Create a resource",
	}

	return cmd
}

// SetOptions is how we get flag values from parent commands. The root command
// has a debug flag. We don't want to copy and paste that flag to every command
// nor do we want to take that flag and put it in a different package so all the
// other commands can import it. We simply need the value. SetOptions lets us
// get that value.
func (cc *CreateCommand) SetOptions(flags cli.Flags) error {
	cc.Debug = cli.ValueOf[bool](flags, "debug")

	return nil
}

func (CreateCommand) Run() error {
	return cli.ErrPrintHelp
}

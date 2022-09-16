package app

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rdeusser/cli"
)

type CreateCommand struct {
	Test     bool
	Filename cli.Path
	Name     string
	Region   string
}

func (cc *CreateCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "create",
		Desc: "Creates some things, probably",
		Flags: cli.Flags{
			&cli.Flag[bool]{
				Name:      "test",
				Shorthand: "t",
				Desc:      "Create a test server",
				Value:     &cc.Test,
				Required:  false,
			},
			&cli.Flag[cli.Path]{
				Name:      "filename",
				Shorthand: "f",
				Desc:      "Read some stuff from a file",
				Value:     &cc.Filename,
				Required:  true,
			},
		},
		Args: cli.Args{
			&cli.Arg[string]{
				Name:     "name",
				Desc:     "Server name",
				Value:    &cc.Name,
				Required: true,
			},
			&cli.Arg[string]{
				Name:     "region",
				Desc:     "Region to create the server in",
				Value:    &cc.Region,
				Required: true,
			},
		},
	}

	return cmd
}

func (*CreateCommand) PersistentPreRun() error {
	fmt.Println("[create] running create persistent pre-runner")
	return nil
}

func (*CreateCommand) PreRun() error {
	fmt.Println("[create] running create pre-runner")
	return nil
}

func (cc *CreateCommand) Run() error {
	fmt.Println("[create] running create runner")

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)

	fmt.Println("=========================================")
	fmt.Fprintf(writer, "Is test?\t%t\n", cc.Test)
	fmt.Fprintf(writer, "File path\t%s\n", cc.Filename.Path)
	fmt.Fprintf(writer, "File extension\t%s\n", cc.Filename.Ext)
	fmt.Fprintf(writer, "Is directory?\t%t\n", cc.Filename.IsDir)
	fmt.Fprintf(writer, "File exists\t%t\n", cc.Filename.Exists)
	fmt.Fprintf(writer, "Server name\t%s\n", cc.Name)
	fmt.Fprintf(writer, "Region name\t%s\n", cc.Region)
	writer.Flush()
	fmt.Println("=========================================")

	return nil
}

func (*CreateCommand) PostRun() error {
	fmt.Println("[create] running create post-runner")
	return nil
}

func (*CreateCommand) PersistentPostRun() error {
	fmt.Println("[create] running create persistent post-runner")
	return nil
}

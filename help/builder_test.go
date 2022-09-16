package help

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli/tablewriter"
)

func TestBuilder(t *testing.T) {
	builder := NewBuilder(
		WithNoColor(),
	)
	indent := 4
	padding := 4
	want := `kubectl controls the Kubernetes cluster manager

USAGE:
    kubectl [command] [flags]

COMMANDS:
    apply     Apply a resource
    create    Create a resource
    delete    Delete a resource
    get       Get a resource

FLAGS:
    -A                 All namespaces
        --debug        Set logging level to debug
    -h, --help         Print help information
    -n, --namespace    Namespace to operate on

Use "kubectl [command] --help" for more information about a command.
`

	builder.Text("kubectl controls the Kubernetes cluster manager")
	builder.Newline()
	builder.Newline()
	builder.Header("USAGE:")
	builder.Newline()
	builder.Text(builder.WithIndent("kubectl [command] [flags]", 4))

	commands := tablewriter.NewWriter()

	builder.Newline()
	builder.Newline()
	builder.Header("COMMANDS:")
	builder.Newline()

	commands.AddLine(
		tablewriter.Cell{
			Indent:  indent,
			Padding: padding,
			Text:    "apply",
		},
		tablewriter.Cell{
			Text: "Apply a resource",
		},
	)
	commands.AddLine(
		tablewriter.Cell{
			Indent:  indent,
			Padding: padding,
			Text:    "create",
		},
		tablewriter.Cell{
			Text: "Create a resource",
		},
	)
	commands.AddLine(
		tablewriter.Cell{
			Indent:  indent,
			Padding: padding,
			Text:    "delete",
		},
		tablewriter.Cell{
			Text: "Delete a resource",
		},
	)
	commands.AddLine(
		tablewriter.Cell{
			Indent:  indent,
			Padding: padding,
			Text:    "get",
		},
		tablewriter.Cell{
			Text: "Get a resource",
		},
	)

	builder.Table(commands)

	flags := tablewriter.NewWriter()

	builder.Newline()
	builder.Newline()
	builder.Header("FLAGS:")
	builder.Newline()

	flags.AddLine(
		tablewriter.Cell{
			Indent: indent,
			Text:   "-A",
		},
		tablewriter.Cell{
			Padding: padding,
			Text:    "",
		},
		tablewriter.Cell{
			Text: "All namespaces",
		},
	)
	flags.AddLine(
		tablewriter.Cell{
			Indent: indent,
			Text:   "",
		},
		tablewriter.Cell{
			Padding: padding,
			Text:    "--debug",
		},
		tablewriter.Cell{
			Text: "Set logging level to debug",
		},
	)
	flags.AddLine(
		tablewriter.Cell{
			Indent: indent,
			Text:   "-h",
			Suffix: ", ",
		},
		tablewriter.Cell{
			Padding: padding,
			Text:    "--help",
		},
		tablewriter.Cell{
			Text: "Print help information",
		},
	)
	flags.AddLine(
		tablewriter.Cell{
			Indent: indent,
			Text:   "-n",
			Suffix: ", ",
		},
		tablewriter.Cell{
			Padding: padding,
			Text:    "--namespace",
		},
		tablewriter.Cell{
			Text: "Namespace to operate on",
		},
	)

	builder.Table(flags)
	builder.Newline()
	builder.Newline()
	builder.Text("Use \"kubectl [command] --help\" for more information about a command.")
	builder.Newline()

	assert.Equal(t, want, builder.String())
}

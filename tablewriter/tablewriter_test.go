package tablewriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	writer := NewWriter()
	writer.AddLine(
		Cell{
			Text: "-A",
		},
		Cell{
			Text:    "",
			Padding: 4,
		},
		Cell{
			Text: "All the things",
		},
	)
	writer.AddLine(
		Cell{
			Text: "",
		},
		Cell{
			Text:    "--debug",
			Padding: 4,
		},
		Cell{
			Text: "Set logging level to debug",
		},
	)
	writer.AddLine(
		Cell{
			Text:   "-h",
			Suffix: ", ",
		},
		Cell{
			Text:    "--help",
			Padding: 4,
		},
		Cell{
			Text: "Print help information",
		},
	)

	want := `-A             All the things
    --debug    Set logging level to debug
-h, --help     Print help information`

	s, err := writer.Render()
	assert.NoError(t, err)
	assert.Equal(t, want, s)
}

func TestRenderWithIndent(t *testing.T) {
	writer := NewWriter()
	writer.AddLine(
		Cell{
			Indent: 4,
			Text:   "-A",
		},
		Cell{
			Text:    "",
			Padding: 4,
		},
		Cell{
			Text: "All the things",
		},
	)
	writer.AddLine(
		Cell{
			Indent: 4,
			Text:   "",
		},
		Cell{
			Text:    "--debug",
			Padding: 4,
		},
		Cell{
			Text: "Set logging level to debug",
		},
	)
	writer.AddLine(
		Cell{
			Indent: 4,
			Text:   "-h",
			Suffix: ", ",
		},
		Cell{
			Text:    "--help",
			Padding: 4,
		},
		Cell{
			Text: "Print help information",
		},
	)

	want := `    -A             All the things
        --debug    Set logging level to debug
    -h, --help     Print help information`

	s, err := writer.Render()
	assert.NoError(t, err)
	assert.Equal(t, want, s)
}

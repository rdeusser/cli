package tablewriter

import (
	"strings"
	"unicode/utf8"
)

const (
	escapeStartRune = rune(27) // \x1b
	escapeStopRune  = 'm'
)

// Writer represents a grid of rows and columns.
type Writer struct {
	sb    strings.Builder
	lines [][]Cell
	err   error
}

// NewWriter returns an initialized writer.
func NewWriter() *Writer {
	writer := &Writer{
		sb:    strings.Builder{},
		lines: make([][]Cell, 0),
	}

	return writer
}

// AddLine adds a line to the table with the provided cells.
func (w *Writer) AddLine(cells ...Cell) {
	w.lines = append(w.lines, make([]Cell, len(cells)))

	i := len(w.lines) - 1
	for j, cell := range cells {
		text := indent(cell.Indent) + cell.String()

		w.lines[i][j] = Cell{
			Indent: cell.Indent,
			Text:   text,
			Suffix: cell.Suffix,
			Align:  cell.Align,
			size:   len(strings.TrimSpace(text)),
			width:  countRunes(text),
		}
	}
}

// MustRender panics if rendering returns an error.
func (w *Writer) MustRender() string {
	s, err := w.Render()
	if err != nil {
		panic(err)
	}

	return s
}

// Render writes each cell with padding and returns the table as a string.
func (w *Writer) Render() (string, error) {
	if w.err != nil {
		return "", w.err
	}

	for i, line := range w.lines {
		for j, cell := range line {
			padding := w.padding(cell, j)

			switch cell.Align {
			case AlignLeft:
				w.write(cell.Text)
				if !isLast(line, j) {
					w.writePadding(padding)
				}
			case AlignCenter:
				// TODO(rdeusser)
			case AlignRight:
				w.writePadding(padding)
				w.write(cell.Text)
			}
		}

		if !isLast(w.lines, i) {
			w.write("\n")
		}
	}

	return w.sb.String(), nil
}

// padding returns the number of spaces needed to pad a cell.
func (w *Writer) padding(cell Cell, line int) int {
	maxLen := maxColumnWidth(w.lines, line)
	return maxLen - cell.width
}

// writePadding takes a number as input and writes the number of spaces to an
// underlying buffer.
func (w *Writer) writePadding(padding int) {
	if padding < 0 {
		padding = 0
	}

	w.write(strings.Repeat(" ", padding))
}

// write writes the provided string to an underlying buffer.
func (w *Writer) write(s string) {
	w.sb.WriteString(s)
}

// maxColumnWidth takes a slice of slice of cells along with a column and
// returns the largest cell's width.
func maxColumnWidth(lines [][]Cell, column int) int {
	maxLen := 0

	for i := 0; i < len(lines); i++ {
		for _, line := range lines {
			text := stripEscapeSequences(line[column].Text)

			if maxLen == 0 {
				maxLen = len(text)
			}

			if len(text) > maxLen {
				maxLen = len(text)
			}
		}
	}

	return maxLen
}

// countRunes returns the total rune width minus ANSI-style escape sequences.
func countRunes(s string) int {
	text := stripEscapeSequences(s)
	return utf8.RuneCountInString(text)
}

// stripEscapeSequences strips ANSI-style escape sequences from the string. This
// is used to get the width of a string for determining the widest cell in a
// column.
func stripEscapeSequences(s string) string {
	var sb strings.Builder

	isEscapeSequence := false

	for _, r := range s {
		switch {
		case r == escapeStartRune:
			isEscapeSequence = true
		case isEscapeSequence:
			if r == escapeStopRune {
				isEscapeSequence = false
			}
		default:
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// indent is just a helper function that returns repeated spaces specified by n.
func indent(n int) string {
	if n < 0 {
		n = 0
	}

	return strings.Repeat(" ", n)
}

type slicer interface {
	[]Cell | [][]Cell
}

func isLast[T slicer](slice T, index int) bool {
	return index >= len(slice)-1
}

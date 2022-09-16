package tablewriter

import "strings"

// Cell represents a single Cell in a table.
type Cell struct {
	// Indent is the level of indent to give each line. Useful for displaying
	// tables under headers (like in help text).
	Indent int

	// Padding is the padding to add to the right side of the text in a cell.
	Padding int

	// Text is the text in the cell.
	Text string

	// Suffix is a string that will be appended to the text of the cell.
	Suffix string

	// Align is where the text will be aligned in the cell.
	Align Align

	size  int // cell size in bytes (no padding)
	width int // cell width in runes (with padding)
}

// String returns a cell as a string with the text, a suffix, and the padding.
func (c Cell) String() string {
	var sb strings.Builder

	if c.Text == "" {
		return ""
	}

	sb.WriteString(c.Text)
	sb.WriteString(c.Suffix)
	sb.WriteString(strings.Repeat(" ", c.Padding))

	return sb.String()
}

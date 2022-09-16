package help

import (
	"fmt"
	"strings"

	"github.com/rdeusser/cli/internal/termenv"
	"github.com/rdeusser/cli/tablewriter"
)

// Builder lets you build help information for cli's and services.
type Builder struct {
	colorize bool
	sb       strings.Builder
}

// Option lets you provide options to Builder.
type Option func(*Builder)

// WithNoColor disables colorization of output.
func WithNoColor() Option {
	return func(b *Builder) {
		b.colorize = false
	}
}

// NewBuilder initializes a new builder.
func NewBuilder(options ...Option) *Builder {
	b := &Builder{
		colorize: true,
		sb:       strings.Builder{},
	}

	for _, option := range options {
		option(b)
	}

	return b
}

// Header writes yellow text to the builder.
func (b *Builder) Header(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	b.WriteString(b.Yellow(s))
}

// Text writes plan text to the builder.
func (b *Builder) Text(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	b.WriteString(s)
}

// WithIndent returns a string with spaces equal to that of indent.
func (*Builder) WithIndent(s string, indent int) string {
	return strings.Repeat(" ", indent) + s
}

// Table writes a table writer to the builder.
func (b *Builder) Table(table *tablewriter.Writer) {
	b.WriteString(table.MustRender())
}

// Newline writes a newline to the builder.
func (b *Builder) Newline() {
	b.WriteString("\n")
}

// Yellow returns text colored as yellow. Respects the WithNoColor option.
func (b *Builder) Yellow(format string, a ...any) string {
	if len(a) == 1 && a[0] == "" {
		return ""
	}

	if b.colorize {
		return termenv.Yellow(format, a...)
	}

	return fmt.Sprintf(format, a...)
}

// Green returns text colored as green. Respects the WithNoColor option.
func (b *Builder) Green(format string, a ...any) string {
	if len(a) == 1 && a[0] == "" {
		return ""
	}

	if b.colorize {
		return termenv.Green(format, a...)
	}

	return fmt.Sprintf(format, a...)
}

// Write implements io.Writer.
func (b *Builder) Write(p []byte) (n int, err error) {
	return b.sb.Write(p)
}

// WriteString writes a string to the builder.
func (b *Builder) WriteString(s string) (n int, err error) {
	return b.sb.WriteString(s)
}

// String returns all the text provided to the builder as a string.
func (b *Builder) String() string {
	return b.sb.String()
}

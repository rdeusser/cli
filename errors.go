package cli

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/termenv"
)

var (
	// ErrPrintHelp indicates that we should show the help output.
	ErrPrintHelp = errors.New("help")

	// ErrInvalidShorthand indicates that too many characters were assigned
	// to the shorthand portion of a flag.
	ErrInvalidShorthand = errors.New("shorthand must be a single character")
)

// ErrFlagAlreadyDefined is when you attempt to add a flag that has already been
// added.
type ErrFlagAlreadyDefined struct {
	option FlagOptionGetter
}

// Error returns an error string of what flag of some type was already defined in
// the command.
func (e ErrFlagAlreadyDefined) Error() string {
	if e.option == nil {
		return termenv.Colorize(termenv.ColorRed, "flag already defined, but is nil")
	}

	option := e.option.Option()

	return fmt.Sprintf("%s %v", termenv.Colorize(termenv.ColorRed, "flag (type %s) already defined:", option.Type), option.Name)
}

// ErrFlagNotDefined is when you attempt to pass a flag to the command that was
// not added to the command itself.
type ErrFlagNotDefined struct {
	flag string
}

// Error returns an error string of which flag was passed to the command, but
// not added to it before the parsing stage.
func (e ErrFlagNotDefined) Error() string {
	return fmt.Sprintf("* %s %v\n", termenv.Colorize(termenv.ColorRed, "flag provided but not defined:"), e.flag)
}

// ErrArgAlreadyDefined is when you attempt to add an argument to a command that
// was already added.
//
// This is mainly used if you attempt to pass an argument to a command and they
// have the same position or the same name.
type ErrArgAlreadyDefined struct {
	option ArgOptionGetter
}

// Error returns an error string of what argument you attempted to construct
// within the command that was already defined.
func (e ErrArgAlreadyDefined) Error() string {
	if e.option == nil {
		return termenv.Colorize(termenv.ColorRed, "arg already defined, but is nil")
	}

	option := e.option.Option()

	return fmt.Sprintf("%s %v", termenv.Colorize(termenv.ColorRed, "arg (position %d, type %s) already defined:", option.Position, option.Type), option.Name)
}

// ErrArgNotDefined is when you attempt to pass an argument to a command that
// was not added to the command itself during the parsing stage.
type ErrArgNotDefined struct {
	arg      string
	position int
}

// Error returns an error string of what argument you passed to the command, but
// did not add to the command itself.
func (e ErrArgNotDefined) Error() string {
	return fmt.Sprintf("* %s %v\n", termenv.Colorize(termenv.ColorRed, "arg (position %d) provided but not defined:", e.position), e.arg)
}

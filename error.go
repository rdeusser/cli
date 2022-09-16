package cli

import (
	"strings"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/termenv"
)

var (
	// ErrPrintHelp indicates that we should show the help output.
	ErrPrintHelp = errors.New("help")

	// ErrEnvVarMustHaveName indicates that you provided a flag or arg with an
	// environment variable that doesn't have a name.
	ErrEnvVarMustHaveName = errors.New("environment variable must have a name")

	// ErrInvalidShorthand indicates that too many characters were assigned
	// to the shorthand portion of a flag.
	ErrInvalidShorthand = errors.New("shorthand must be a single letter")

	// ErrFlagSliceMustHaveSeparator indicates that a separator was not set on a
	// flag or is "".
	ErrFlagSliceMustHaveSeparator = errors.New("flag must have a separator if value is a slice")

	// ErrMustHaveParent indicates that an OptionSetter was given to a command
	// that doesn't have a parent (e.g. root command).
	ErrMustHaveParent = errors.New("command must have parent in order to use SetOptions")
)

// ErrFlagAlreadyDefined is when you attempt to add a flag that has already been
// added.
type ErrFlagAlreadyDefined struct {
	Name      string
	Shorthand string
}

// Error returns an error string of what flag of some type was already defined in
// the command.
func (e ErrFlagAlreadyDefined) Error() string {
	if e.Shorthand != "" {
		return termenv.Red("-%s, --%s already defined", e.Shorthand, e.Name)
	}

	return termenv.Red("--%s already defined", e.Name)
}

// ErrArgAlreadyDefined is when you attempt to add an argument to a command that
// was already added.
//
// This is mainly used if you attempt to pass an argument to a command and they
// have the same position or the same name.
type ErrArgAlreadyDefined struct {
	Name string
}

// Error returns an error string of what argument you attempted to construct
// within the command that was already defined.
func (e ErrArgAlreadyDefined) Error() string {
	return termenv.Red("<%s> already defined", e.Name)
}

// ErrFlagRequired is an error describing a flag that is required.
type ErrFlagRequired struct {
	Name      string
	Shorthand string
}

// Error returns an error string when you don't pass a required flag on the
// command line.
func (e ErrFlagRequired) Error() string {
	if e.Shorthand != "" {
		return termenv.Red("-%s, --%s is required", e.Shorthand, e.Name)
	}

	return termenv.Red("--%s is required", e.Name)
}

// ErrArgRequired is an error describing an argument that is required.
type ErrArgRequired struct {
	Name string
}

// Error returns an error string when you don't pass a required argument on the
// command line.
func (e ErrArgRequired) Error() string {
	return termenv.Red("<%s> is required", e.Name)
}

// ErrUnknown is an error describing an argument or flag that wasn't defined.
type ErrUnknown struct {
	Input    string
	Arg      string
	StartPos int
	EndPos   int
}

// Error returns an error string for describing an argument or flag that wasn't defined.
func (e ErrUnknown) Error() string {
	var lb lineBuilder

	lb.Write(termenv.Red("warning: "))
	lb.Write(termenv.BrightWhite("unknown argument '%s'", e.Arg))
	lb.NewLine()
	lb.Write("\t")
	lb.Write(e.Input)

	idx := lastIndex(lb.CurrentLine(), e.Arg)

	lb.NewLine()
	lb.Write(columnToSpace(idx))

	count := e.EndPos - e.StartPos - 1
	if count < 0 {
		count = 0
	}

	lb.Write(strings.Repeat(termenv.Red("^"), count))
	lb.Flush()

	return lb.String()
}

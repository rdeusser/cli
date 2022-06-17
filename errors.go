package cli

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrPrintHelp        = errors.New("help")
	ErrInvalidShorthand = errors.New("shorthand must be a single letter/number")
)

type ErrFlagAlreadyDefined struct {
	option FlagOptionGetter
}

func (e ErrFlagAlreadyDefined) Error() string {
	if e.option == nil {
		return colorize(ColorRed, "flag already defined, but is nil")
	}

	option := e.option.Option()

	return fmt.Sprintf("%s %v", colorize(ColorRed, "flag (type %s) already defined:", option.Type()), option.Name)
}

type ErrFlagNotDefined struct {
	flag string
}

func (e ErrFlagNotDefined) Error() string {
	return fmt.Sprintf("* %s %v\n", colorize(ColorRed, "flag provided but not defined:"), e.flag)
}

type ErrArgAlreadyDefined struct {
	option ArgOptionGetter
}

func (e ErrArgAlreadyDefined) Error() string {
	if e.option == nil {
		return colorize(ColorRed, "arg already defined, but is nil")
	}

	option := e.option.Option()

	return fmt.Sprintf("%s %v", colorize(ColorRed, "arg (type %s) already defined:", option.Type()), option.Name)
}

type ErrArgNotDefined struct {
	arg string
}

func (e ErrArgNotDefined) Error() string {
	return fmt.Sprintf("* %s %v\n", colorize(ColorRed, "arg provided but not defined:"), e.arg)
}

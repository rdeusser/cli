package cli

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	PrintHelp           = errors.New("help")
	ErrInvalidShorthand = errors.New("shorthand must be a single letter/number")
)

type ErrOptionAlreadyDefined struct {
	opt OptionGetter
}

func (e ErrOptionAlreadyDefined) Error() string {
	if e.opt == nil {
		return colorize(ColorRed, "option already defined, but is nil")
	}

	opt, err := e.opt.Option()
	if err != nil {
		return errors.Wrapf(err, "%s %v", colorize(ColorRed, "option (type %s) already defined:", opt.Type()), opt.Name).Error()
	}

	return fmt.Sprintf("%s %v", colorize(ColorRed, "option (type %s) already defined:", opt.Type()), opt.Name)
}

type ErrOptionNotDefined struct {
	arg string
}

func (e ErrOptionNotDefined) Error() string {
	return fmt.Sprintf("* %s %v\n", colorize(ColorRed, "option provided but not defined:"), e.arg)
}

package cli

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	PrintHelp           = errors.New("help")
	ErrInvalidShorthand = errors.New("shorthand must be a single letter/number")
)

type ErrOptionNotDefined struct {
	opt OptionGetter
	arg string
}

func (e ErrOptionNotDefined) Error() string {
	if e.opt == (*Option)(nil) {
		return fmt.Sprintf("%s %v\n", bad("* option provided but not defined:"), e.arg)
	}

	opt, err := e.opt.GetOption()
	if err != nil {
		return errors.Wrapf(err, "%s %v\n", bad("* option (type %s) provided but not defined:", opt.GetType()), e.arg).Error()
	}

	return fmt.Sprintf("%s %v\n", bad("* option (type %s) provided but not defined:", opt.GetType()), e.arg)
}

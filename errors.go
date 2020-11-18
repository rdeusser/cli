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
	if e.opt == (*Option)(nil) {
		return fmt.Sprintf("%s", bad("option already defined, but is nil"))
	}

	opt, err := e.opt.GetOption()
	if err != nil {
		return errors.Wrapf(err, "%s %v", bad("option (type %s) already defined:", opt.GetType()), opt.Name).Error()
	}

	return fmt.Sprintf("%s %v", bad("option (type %s) already defined:", opt.GetType()), opt.Name)
}

type ErrOptionNotDefined struct {
	arg string
}

func (e ErrOptionNotDefined) Error() string {
	return fmt.Sprintf("* %s %v\n", bad("option provided but not defined:"), e.arg)
}

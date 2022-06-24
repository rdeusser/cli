package cli

import (
	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/humanize"
	"github.com/rdeusser/cli/internal/types"
)

// BoolArg is a bool argument.
type BoolArg struct {
	Bind     *bool
	Name     string
	Desc     string
	Position int
	Required bool

	option ArgOption
	value  *types.Bool
}

// String returns a string-formatted bool value.
func (a *BoolArg) String() string {
	return a.value.String()
}

// Set sets the bool argument's value.
func (a *BoolArg) Set(s string) error {
	if err := a.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a bool value for the %s argument", s, humanize.Ordinal(a.Position+1))
	}

	a.option.HasBeenSet = true
	return nil
}

// Get gets the value of the bool argument.
func (a *BoolArg) Get() bool {
	return *a.Bind
}

// Type returns the type of the argument.
func (a *BoolArg) Type() types.Type {
	return types.BoolType
}

// Option returns the option for the arg.
func (a *BoolArg) Option() ArgOption {
	return a.option
}

// Init initializes the default (or already set) options for the argument. Most
// notably, it doesn't indicate that the argument has actually been set
// yet. That's the job of the parser.
func (a *BoolArg) Init() error {
	if a.value == nil {
		a.value = types.NewBool(a.Bind, false)
	}

	a.option = ArgOption{
		Bind:       a.Bind,
		Name:       a.Name,
		Desc:       a.Desc,
		Position:   a.Position,
		Required:   a.Required,
		Type:       a.Type(),
		HasBeenSet: false,
	}

	return nil
}

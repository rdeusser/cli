package cli

import (
	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/humanize"
	"github.com/rdeusser/cli/internal/types"
)

// IntArg is a string argument.
type IntArg struct {
	Bind     *int
	Name     string
	Desc     string
	Position int
	Required bool

	option ArgOption
	value  *types.Int
}

// String returns a string-formatted string value.
func (a *IntArg) String() string {
	return a.value.String()
}

// Set sets the string argument's value.
func (a *IntArg) Set(s string) error {
	if a.value == nil {
		a.value = types.NewInt(a.Bind, 0)
	}

	if err := a.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a string value for the %s argument", s, humanize.Ordinal(a.Position))
	}

	a.option.hasBeenSet = true
	return nil
}

// Get gets the value of the string argument.
func (a *IntArg) Get() int {
	return *a.Bind
}

// Type returns the type of the argument.
func (a *IntArg) Type() types.Type {
	return types.IntType
}

// Option returns the option for the arg.
func (a *IntArg) Option() ArgOption {
	return a.option
}

// Apply applies the default (or already set) options for the argument. Most
// notably, it doesn't indicate that the argument has actually been set
// yet. That's the job of the parser.
func (a *IntArg) Apply() error {
	if a.value == nil {
		a.value = types.NewInt(a.Bind, 0)
	}

	a.option = ArgOption{
		Bind:     a.Bind,
		Name:     a.Name,
		Desc:     a.Desc,
		Position: a.Position,
		Required: a.Required,

		typ:        a.Type(),
		hasBeenSet: false,
	}

	return nil
}

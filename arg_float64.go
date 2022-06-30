package cli

import (
	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/humanize"
	"github.com/rdeusser/cli/internal/types"
)

// Float64Arg is a float64 argument.
type Float64Arg struct {
	Bind     *float64
	Name     string
	Desc     string
	Position int
	Required bool

	option ArgOption
	value  *types.Float64
}

// String returns a string-formatted float64 value.
func (a *Float64Arg) String() string {
	return a.value.String()
}

// Set sets the float64 argument's value.
func (a *Float64Arg) Set(s string) error {
	if err := a.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a float64 value for the %s argument", s, humanize.Ordinal(a.Position+1))
	}

	a.option.HasBeenSet = true
	return nil
}

// Get gets the value of the float64 argument.
func (a *Float64Arg) Get() float64 {
	return *a.Bind
}

// Type returns the type of the argument.
func (a *Float64Arg) Type() types.Type {
	return types.Float64Type
}

// Option returns the option for the arg.
func (a *Float64Arg) Option() ArgOption {
	return a.option
}

// Init initializes the default (or already set) options for the argument. Most
// notably, it doesn't indicate that the argument has actually been set
// yet. That's the job of the parser.
func (a *Float64Arg) Init() error {
	if a.value == nil {
		a.value = types.NewFloat64(a.Bind, 0.0)
	}

	a.option = ArgOption{
		Bind:     a.Bind,
		Name:     a.Name,
		Desc:     a.Desc,
		Position: a.Position,
		Required: a.Required,
		Type:     a.Type(),
	}

	return nil
}

package cli

import (
	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/humanize"
	"github.com/rdeusser/cli/internal/types"
)

// StringArg is a string argument.
type StringArg struct {
	Bind     *string
	Name     string
	Desc     string
	Position int
	Required bool

	option ArgOption
	value  *types.String
}

// String returns a string-formatted string value.
func (a *StringArg) String() string {
	return a.value.String()
}

// Set sets the string argument's value.
func (a *StringArg) Set(s string) error {
	if err := a.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a string value for the %s argument", s, humanize.Ordinal(a.Position+1))
	}

	a.option.HasBeenSet = true
	return nil
}

// Get gets the value of the string argument.
func (a *StringArg) Get() string {
	return *a.Bind
}

// Type returns the type of the argument.
func (a *StringArg) Type() types.Type {
	return types.StringType
}

// Option returns the option for the arg.
func (a *StringArg) Option() ArgOption {
	return a.option
}

// Init initializes the default (or already set) options for the argument. Most
// notably, it doesn't indicate that the argument has actually been set
// yet. That's the job of the parser.
func (a *StringArg) Init() error {
	if a.value == nil {
		a.value = types.NewString(a.Bind, "")
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

package cli

import (
	"time"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/humanize"
	"github.com/rdeusser/cli/internal/types"
)

// DurationArg is a time.Duration argument.
type DurationArg struct {
	Bind     *time.Duration
	Name     string
	Desc     string
	Position int
	Required bool

	option ArgOption
	value  *types.Duration
}

// String returns a string-formatted time.Duration value.
func (a *DurationArg) String() string {
	return a.value.String()
}

// Set sets the time.Duration argument's value.
func (a *DurationArg) Set(s string) error {
	if err := a.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a time.Duration value for the %s argument", s, humanize.Ordinal(a.Position+1))
	}

	a.option.HasBeenSet = true
	return nil
}

// Get gets the value of the time.Duration argument.
func (a *DurationArg) Get() time.Duration {
	return *a.Bind
}

// Type returns the type of the argument.
func (a *DurationArg) Type() types.Type {
	return types.DurationType
}

// Option returns the option for the arg.
func (a *DurationArg) Option() ArgOption {
	return a.option
}

// Init initializes the default (or already set) options for the argument. Most
// notably, it doesn't indicate that the argument has actually been set
// yet. That's the job of the parser.
func (a *DurationArg) Init() error {
	if a.value == nil {
		a.value = types.NewDuration(a.Bind, time.Duration(0))
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

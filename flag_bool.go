package cli

import (
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/types"
)

// BoolFlag is a bool flag.
type BoolFlag struct {
	Bind      *bool
	Name      string
	Shorthand string
	Desc      string
	Default   bool
	EnvVar    string
	Required  bool

	option FlagOption
	value  *types.Bool
}

// String returns a string-formatted bool value.
func (f *BoolFlag) String() string {
	return f.value.String()
}

// Set sets the bool flag's value.
func (f *BoolFlag) Set(s string) error {
	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		if err := f.value.Set(v); err != nil {
			return errors.Wrapf(err, "setting %s as a bool value for the %s flag", v, f.Name)
		}
	}

	if err := f.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a bool value for the %s flag", s, f.Name)
	}

	f.option.HasBeenSet = true
	return nil
}

// Get gets the value of the bool flag.
func (f *BoolFlag) Get() bool {
	return f.value.Get()
}

// Type returns the type of the flag.
func (f *BoolFlag) Type() types.Type {
	return types.BoolType
}

// Option returns the option for the flag.
func (f *BoolFlag) Option() FlagOption {
	return f.option
}

// Init initializes the default (or already set) options for the flag. Most
// notably, it doesn't indicate that the flag has actually been set yet. That's
// the job of the parser.
func (f *BoolFlag) Init() error {
	if f.value == nil {
		f.value = types.NewBool(f.Bind, f.Default)
	}

	f.option = FlagOption{
		Bind:      f.Bind,
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Default:   f.value.String(),
		Required:  f.Required,
		Type:      f.Type(),
	}

	return nil
}

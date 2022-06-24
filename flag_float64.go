package cli

import (
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/types"
)

// Float64Flag is a float64 flag.
type Float64Flag struct {
	Bind      *float64
	Name      string
	Shorthand string
	Desc      string
	Default   float64
	EnvVar    string
	Required  bool

	option FlagOption
	value  *types.Float64
}

// String returns a string-formatted float64 value.
func (f *Float64Flag) String() string {
	return f.value.String()
}

// Set sets the float64 flag's value.
func (f *Float64Flag) Set(s string) error {
	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		if err := f.value.Set(v); err != nil {
			return errors.Wrapf(err, "setting %s as a float64 value for the %s flag", v, f.Name)
		}
	}

	if err := f.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a float64 value for the %s flag", s, f.Name)
	}

	f.option.HasBeenSet = true
	return nil
}

// Get gets the value of the float64 flag.
func (f *Float64Flag) Get() float64 {
	return f.value.Get()
}

// Type returns the type of the flag.
func (f *Float64Flag) Type() types.Type {
	return types.Float64Type
}

// Option returns the option for the flag.
func (f *Float64Flag) Option() FlagOption {
	return f.option
}

// Init initializes the default (or already set) options for the flag. Most
// notably, it doesn't indicate that the flag has actually been set yet. That's
// the job of the parser.
func (f *Float64Flag) Init() error {
	if f.value == nil {
		f.value = types.NewFloat64(f.Bind, f.Default)
	}

	f.option = FlagOption{
		Bind:       f.Bind,
		Name:       f.Name,
		Shorthand:  f.Shorthand,
		Desc:       f.Desc,
		EnvVar:     f.EnvVar,
		Default:    f.value.String(),
		Required:   f.Required,
		Type:       f.Type(),
		HasBeenSet: false,
	}

	return nil
}

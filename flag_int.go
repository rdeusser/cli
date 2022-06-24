package cli

import (
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/types"
)

// IntFlag is a int flag.
type IntFlag struct {
	Bind      *int
	Name      string
	Shorthand string
	Desc      string
	Default   int
	EnvVar    string
	Required  bool

	option FlagOption
	value  *types.Int
}

// String returns a string-formatted int value.
func (f *IntFlag) String() string {
	return f.value.String()
}

// Set sets the int flag's value.
func (f *IntFlag) Set(s string) error {
	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		if err := f.value.Set(v); err != nil {
			return errors.Wrapf(err, "setting %s as a int value for the %s flag", v, f.Name)
		}
	}

	f.option.HasBeenSet = true
	return nil
}

// Get gets the value of the int flag.
func (f *IntFlag) Get() int {
	return f.value.Get()
}

// Type returns the type of the flag.
func (f *IntFlag) Type() types.Type {
	return types.IntType
}

// Option returns the option for the flag.
func (f *IntFlag) Option() FlagOption {
	return f.option
}

// Init initializes the default (or already set) options for the flag. Most
// notably, it doesn't indicate that the flag has actually been set yet. That's
// the job of the parser.
func (f *IntFlag) Init() error {
	if f.value == nil {
		f.value = types.NewInt(f.Bind, f.Default)
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

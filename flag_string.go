package cli

import (
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/types"
)

// StringFlag is a string flag.
type StringFlag struct {
	Bind      *string
	Name      string
	Shorthand string
	Desc      string
	Default   string
	EnvVar    string
	Required  bool

	option FlagOption
	value  *types.String
}

// String returns a string-formatted string value.
func (f *StringFlag) String() string {
	return f.value.String()
}

// Set sets the string flag's value.
func (f *StringFlag) Set(s string) error {
	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		if err := f.value.Set(v); err != nil {
			return errors.Wrapf(err, "setting %s as a string value for the %s flag", v, f.Name)
		}
	}

	f.option.HasBeenSet = true
	return nil
}

// Get gets the value of the string flag.
func (f *StringFlag) Get() string {
	return f.value.Get()
}

// Type returns the type of the flag.
func (f *StringFlag) Type() types.Type {
	return types.StringType
}

// Option returns the option for the flag.
func (f *StringFlag) Option() FlagOption {
	return f.option
}

// Init initializes the default (or already set) options for the flag. Most
// notably, it doesn't indicate that the flag has actually been set yet. That's
// the job of the parser.
func (f *StringFlag) Init() error {
	if f.value == nil {
		f.value = types.NewString(f.Bind, f.Default)
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

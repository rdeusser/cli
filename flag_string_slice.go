package cli

import (
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/types"
)

// StringSliceFlag is a []string flag.
type StringSliceFlag struct {
	Bind      *[]string
	Name      string
	Shorthand string
	Desc      string
	Default   []string
	EnvVar    string
	Separator string
	Required  bool

	option FlagOption
	value  *types.StringSlice
}

// String returns a string-formatted []string value.
func (f *StringSliceFlag) String() string {
	return f.value.String()
}

// Set sets the []string flag's value.
func (f *StringSliceFlag) Set(s string) error {
	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		if err := f.value.Set(v); err != nil {
			return errors.Wrapf(err, "setting %s as a []string value for the %s flag", v, f.Name)
		}
	}

	if err := f.value.Set(s); err != nil {
		return errors.Wrapf(err, "setting %s as a []string value for the %s flag", s, f.Name)
	}

	f.option.HasBeenSet = true
	return nil
}

// Get gets the value of the []string flag.
func (f *StringSliceFlag) Get() []string {
	return f.value.Get()
}

// Type returns the type of the flag.
func (f *StringSliceFlag) Type() types.Type {
	return types.StringSliceType
}

// Option returns the option for the flag.
func (f *StringSliceFlag) Option() FlagOption {
	return f.option
}

// Init initializes the default (or already set) options for the flag. Most
// notably, it doesn't indicate that the flag has actually been set yet. That's
// the job of the parser.
func (f *StringSliceFlag) Init() error {
	if f.value == nil {
		f.value = types.NewStringSlice(f.Bind, f.Default)
	}

	f.option = FlagOption{
		Bind:       f.Bind,
		Name:       f.Name,
		Shorthand:  f.Shorthand,
		Desc:       f.Desc,
		Default:    f.value.String(),
		EnvVar:     f.EnvVar,
		Separator:  f.Separator,
		Required:   f.Required,
		Type:       f.Type(),
		HasBeenSet: false,
	}

	return nil
}

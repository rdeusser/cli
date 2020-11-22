package cli

import (
	"flag"
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/values"
)

type StringSliceFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   []string
	Value     flag.Value
	EnvVar    string
	Required  bool
}

func (f *StringSliceFlag) Type() OptionType {
	return Strings
}

func (f *StringSliceFlag) Option() (Option, error) {
	value := values.NewStrings(nil, f.Default)

	if len(f.Shorthand) > 1 {
		return Option{}, ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		if err := value.Set(v); err != nil {
			return Option{}, err
		}
	}

	f.Value = value

	return Option{
		optType: Strings,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   value.String(),
		Required:  f.Required,
	}, nil
}

func (f *StringSliceFlag) String() string {
	return f.Value.String()
}

func (f *StringSliceFlag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *StringSliceFlag) Get() []string {
	return strings.Split(f.Value.String(), " ")
}

package cli

import (
	"flag"
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/values"
)

type StringFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   string
	Value     flag.Value
	EnvVar    string
	Required  bool
}

func (f *StringFlag) Type() OptionType {
	return String
}

func (f *StringFlag) Option() (Option, error) {
	value := values.NewString(nil, f.Default)

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
		optType: String,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   value.String(),
		Required:  f.Required,
	}, nil
}

func (f *StringFlag) String() string {
	return f.Value.String()
}

func (f *StringFlag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *StringFlag) Get() string {
	return f.Value.String()
}

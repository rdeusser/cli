package cli

import (
	"flag"
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/values"
)

type StringFlag struct {
	value flag.Value

	Name      string
	Shorthand string
	Desc      string
	Default   string
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

	f.value = value

	return Option{
		optType: String,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     f.value,
		Default:   f.value.String(),
		Required:  f.Required,
	}, nil
}

func (f *StringFlag) String() string {
	if f.value == nil {
		panic("value of string flag is nil, did you add it to your command?")
	}
	return f.value.String()
}

func (f *StringFlag) Set(s string) error {
	return f.value.Set(s)
}

func (f *StringFlag) Get() string {
	return f.value.String()
}

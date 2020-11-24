package cli

import (
	"flag"
	"os"
	"strings"

	"github.com/rdeusser/cli/internal/values"
)

type StringSliceFlag struct {
	value flag.Value

	Name      string
	Shorthand string
	Desc      string
	Default   []string
	EnvVar    string
	Required  bool
}

func (f *StringSliceFlag) Type() OptionType {
	return StringSlice
}

func (f *StringSliceFlag) Option() (Option, error) {
	value := values.NewStringSlice(nil, f.Default)

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
		optType: StringSlice,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     f.value,
		Default:   f.value.String(),
		Required:  f.Required,
	}, nil
}

func (f *StringSliceFlag) String() string {
	return f.value.String()
}

func (f *StringSliceFlag) Set(s string) error {
	return f.value.Set(s)
}

func (f *StringSliceFlag) Get() []string {
	value := f.value.String()
	value = strings.ReplaceAll(value, "[", "")
	value = strings.ReplaceAll(value, "]", "")
	value = strings.ReplaceAll(value, "\"", "")

	return strings.Split(value, ", ")
}

func (f *StringSliceFlag) Clear() bool {
	f.value = values.NewStringSlice(nil, []string{})
	return len(f.String()) > 0
}

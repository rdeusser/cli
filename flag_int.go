package cli

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/values"
)

type IntFlag struct {
	value flag.Value

	Name      string
	Shorthand string
	Desc      string
	Default   int
	EnvVar    string
	Required  bool
}

func (f *IntFlag) Type() OptionType {
	return Int
}

func (f *IntFlag) Option() (Option, error) {
	value := values.NewInt(nil, f.Default)

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		_, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return Option{}, errors.Wrapf(err, "parsing %q as an int value for flag %s", v, f.Name)
		}

		if err := value.Set(v); err != nil {
			return Option{}, err
		}
	}

	f.value = value

	return Option{
		optType: Int,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     f.value,
		Default:   f.value.String(),
		Required:  f.Required,
	}, nil
}

func (f *IntFlag) String() string {
	return f.value.String()
}

func (f *IntFlag) Set(s string) error {
	return f.value.Set(s)
}

func (f *IntFlag) Get() int {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	i, _ := strconv.ParseInt(f.value.String(), 10, 64)
	return int(i)
}

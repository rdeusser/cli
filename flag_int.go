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
	Name      string
	Shorthand string
	Desc      string
	Default   int
	Value     flag.Value
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

	f.Value = value

	return Option{
		optType: Int,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   value.String(),
		Required:  f.Required,
	}, nil
}

func (f *IntFlag) String() string {
	return f.Value.String()
}

func (f *IntFlag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *IntFlag) Get() int {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	i, _ := strconv.ParseInt(f.Value.String(), 10, 64)
	return int(i)
}

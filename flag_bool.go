package cli

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/values"
)

type BoolFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   bool
	Value     flag.Value
	EnvVar    string
	Required  bool
}

func (f *BoolFlag) Type() OptionType {
	return Bool
}

func (f *BoolFlag) Option() (Option, error) {
	value := values.NewBool(nil, f.Default)

	if len(f.Shorthand) > 1 {
		return Option{}, ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		_, err := strconv.ParseBool(v)
		if err != nil {
			return Option{}, errors.Wrapf(err, "parsing %q as a bool value for flag %s", v, f.Name)
		}

		if err := value.Set(v); err != nil {
			return Option{}, err
		}
	}

	f.Value = value

	return Option{
		optType: Bool,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   value.String(),
		Required:  f.Required,
	}, nil
}

func (f *BoolFlag) String() string {
	return f.Value.String()
}

func (f *BoolFlag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *BoolFlag) Get() bool {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	b, _ := strconv.ParseBool(f.Value.String())
	return b
}

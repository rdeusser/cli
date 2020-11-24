package cli

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/values"
)

type Float64Flag struct {
	value flag.Value

	Name      string
	Shorthand string
	Desc      string
	Default   float64
	EnvVar    string
	Required  bool
}

func (f *Float64Flag) Type() OptionType {
	return Float64
}

func (f *Float64Flag) Option() (Option, error) {
	value := values.NewFloat64(nil, f.Default)

	if len(f.Shorthand) > 1 {
		return Option{}, ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return Option{}, errors.Wrapf(err, "parsing %q as a float64 value for flag %s", v, f.Name)
		}

		if err := value.Set(v); err != nil {
			return Option{}, err
		}
	}

	f.value = value

	return Option{
		optType: Float64,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     f.value,
		Default:   f.value.String(),
		Required:  f.Required,
	}, nil
}

func (f *Float64Flag) String() string {
	return f.value.String()
}

func (f *Float64Flag) Set(s string) error {
	return f.value.Set(s)
}

func (f *Float64Flag) Get() float64 {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	f64, _ := strconv.ParseFloat(f.value.String(), 0)
	return f64
}

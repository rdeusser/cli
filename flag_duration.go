package cli

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/values"
)

type DurationFlag struct {
	value flag.Value

	Name      string
	Shorthand string
	Desc      string
	Default   time.Duration
	EnvVar    string
	Required  bool
}

func (f *DurationFlag) Type() OptionType {
	return Duration
}

func (f *DurationFlag) Option() (Option, error) {
	value := values.NewDuration(nil, f.Default)

	if len(f.Shorthand) > 1 {
		return Option{}, ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		_, err := time.ParseDuration(v)
		if err != nil {
			return Option{}, errors.Wrapf(err, "parsing %q as a time.Duration value for flag %s", v, f.Name)
		}

		if err := value.Set(v); err != nil {
			return Option{}, err
		}
	}

	f.value = value

	return Option{
		optType: Duration,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     f.value,
		Default:   f.value.String(),
		Required:  f.Required,
	}, nil
}

func (f *DurationFlag) String() string {
	if f.value == nil {
		panic("value of duration flag is nil, did you add it to your command?")
	}
	return f.value.String()
}

func (f *DurationFlag) Set(s string) error {
	return f.value.Set(s)
}

func (f *DurationFlag) Get() time.Duration {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	d, _ := time.ParseDuration(f.String())
	return d
}

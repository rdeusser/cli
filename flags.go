package cli

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/values"
)

var (
	HelpFlag = BoolFlag{
		Name:      "help",
		Shorthand: "h",
		Desc:      "show help",
	}

	VersionFlag = BoolFlag{
		Name:      "version",
		Shorthand: "V",
		Desc:      "print the version",
	}
)

type Flag interface {
	flag.Value
	TypeGetter
	OptionGetter
}

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

type Float64Flag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   float64
	Value     flag.Value
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

	f.Value = value

	return Option{
		optType: Float64,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   value.String(),
		Required:  f.Required,
	}, nil
}

func (f *Float64Flag) String() string {
	return f.Value.String()
}

func (f *Float64Flag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *Float64Flag) Get() float64 {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	f64, _ := strconv.ParseFloat(f.Value.String(), 0)
	return f64
}

type DurationFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   time.Duration
	Value     flag.Value
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

	f.Value = value

	return Option{
		optType: Duration,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   f.Default.String(),
		Required:  f.Required,
	}, nil
}

func (f *DurationFlag) String() string {
	return f.Value.String()
}

func (f *DurationFlag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *DurationFlag) Get() time.Duration {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	d, _ := time.ParseDuration(f.Value.String())
	return d
}

type StringsFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   []string
	Value     flag.Value
	EnvVar    string
	Required  bool
}

func (f *StringsFlag) Type() OptionType {
	return Strings
}

func (f *StringsFlag) Option() (Option, error) {
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

func (f *StringsFlag) String() string {
	return f.Value.String()
}

func (f *StringsFlag) Set(s string) error {
	return f.Value.Set(s)
}

func (f *StringsFlag) Get() []string {
	return strings.Split(f.Value.String(), " ")
}

package cli

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/types"
)

type BoolFlag struct {
	option FlagOption

	Name      string
	Shorthand string
	Desc      string
	Default   bool
	EnvVar    string
	Required  bool
}

func (f *BoolFlag) Set(s string) error {
	return f.option.Value.Set(s)
}

func (f *BoolFlag) String() string {
	if f.option.Value == nil {
		panic("value of bool flag is nil, did you add it to your command?")
	}
	return f.option.Value.String()
}

func (f *BoolFlag) Get() bool {
	// By this time, we've already validated the flag so we don't need to do
	// so again.
	b, _ := strconv.ParseBool(f.String())
	return b
}

func (f *BoolFlag) Type() types.Type {
	return types.BoolType
}

func (f *BoolFlag) Option() FlagOption {
	return f.option
}

func (f *BoolFlag) Apply() error {
	value := types.NewBool(nil, f.Default)

	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		_, err := strconv.ParseBool(v)
		if err != nil {
			return errors.Wrapf(err, "parsing %q as a bool value for the %s flag", v, f.Name)
		}

		if err := value.Set(v); err != nil {
			return errors.Wrap(err, "setting %q as a bool value for the %s flag")
		}
	}

	f.option = FlagOption{
		typ: types.BoolType,

		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   value.String(),
		Required:  f.Required,
	}

	return nil
}

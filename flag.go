package cli

import (
	stdflag "flag"

	"github.com/rdeusser/cli/internal/types"
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

// FlagOptionGetter is an interface to indicate that a type provides flagument
// options.
type FlagOptionGetter interface {
	Option() FlagOption
}

// Flags is a collection of flags.
type Flags []Flag

// Lookup looks up a flag by either it's shorthand form or the full name.
func (f *Flags) Lookup(shorthand, name string) Flag {
	for _, flag := range *f {
		option := flag.Option()

		if option.Shorthand == shorthand || option.Name == name {
			return flag
		}
	}

	return nil
}

// Flag is an interface for defining flags.
type Flag interface {
	stdflag.Value
	types.Getter
	FlagOptionGetter

	Apply() error
}

// FlagOption represents all possible underlying flag types.
type FlagOption struct {
	typ types.Type

	Name      string
	Shorthand string
	Desc      string
	Default   string
	Value     stdflag.Value
	EnvVar    string
	Required  bool
}

// Type returns the type of the flag.
func (o FlagOption) Type() types.Type {
	return o.typ
}

// HasBeenSet indicates if the flag has been set.
func (o FlagOption) HasBeenSet() bool {
	return o.Value.String() != ""
}

// SortFlagOptionsByName sorts flags by name.
type SortFlagOptionsByName []FlagOption

func (n SortFlagOptionsByName) Len() int           { return len(n) }
func (n SortFlagOptionsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortFlagOptionsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

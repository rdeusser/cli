package cli

import (
	"github.com/rdeusser/cli/internal/types"
)

var (
	HelpFlag = BoolFlag{
		Name:      "help",
		Shorthand: "h",
		Desc:      "Print help information",
	}

	VersionFlag = BoolFlag{
		Name:      "version",
		Shorthand: "V",
		Desc:      "Print version information",
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

		if option.Name == name {
			return flag
		}

		if option.Shorthand == name {
			return flag
		}
	}

	return nil
}

// Flag is an interface for defining flags.
type Flag interface {
	Value
	types.Getter
	FlagOptionGetter

	Init() error
}

// FlagOption represents all possible underlying flag types.
type FlagOption struct {
	// Bind is a variable or field to set.
	Bind interface{}

	// Name is the name of the flag.
	Name string

	// Shorthand is the short version of the name
	// (e.g. -h inplace of --help).
	Shorthand string

	// Desc is the description of this flag.
	Desc string

	// Default is the default value for this option
	Default string

	// EnvVar is the environment variable to set the flag to if the flag
	// itself wasn't provided.
	EnvVar string

	// Separator is the separator to use when providing multiple arguments
	// to a flag (i.e. a slice).
	Separator string

	// Required indicates whether this flag is required.
	Required bool

	// Type represents the underlying flag type.
	Type types.Type

	// HasBeenSet indicates whether or not the flag was set explicitly.
	//
	// The purpose of this field is to distinguish between a default value
	// and when an flag was explicitly set.
	HasBeenSet bool
}

// String returns the name of the flag.
func (fo FlagOption) String() string {
	return fo.Name
}

// SortFlagOptionsByName sorts flags by name.
type SortFlagOptionsByName []FlagOption

func (n SortFlagOptionsByName) Len() int           { return len(n) }
func (n SortFlagOptionsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortFlagOptionsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

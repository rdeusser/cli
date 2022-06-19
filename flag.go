package cli

import (
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
	Value
	types.Getter
	FlagOptionGetter

	Apply() error
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

	// Required indicates whether this flag is required.
	Required bool

	// typ represents the underlying flag type.
	typ types.Type

	// hasBeenSet indicates whether or not the flag was set explicitly.
	//
	// The purpose of this field is to distinguish between a default value
	// and when an flag was explicitly set.
	hasBeenSet bool
}

// Type returns the type of the flag.
func (o FlagOption) Type() types.Type {
	return o.typ
}

// HasBeenSet indicates if the flag has been set.
func (o FlagOption) HasBeenSet() bool {
	return o.hasBeenSet
}

// SortFlagOptionsByName sorts flags by name.
type SortFlagOptionsByName []FlagOption

func (n SortFlagOptionsByName) Len() int           { return len(n) }
func (n SortFlagOptionsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortFlagOptionsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

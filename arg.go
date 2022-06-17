package cli

import (
	"flag"

	"github.com/rdeusser/cli/internal/types"
)

// ArgOptionGetter is an interface to indicate that a type provides argument
// options.
type ArgOptionGetter interface {
	Option() ArgOption
}

// Args is a collection of arguments.
type Args []Arg

// Lookup looks up an argument by name.
func (a *Args) Lookup(name string) Arg {
	for _, arg := range *a {
		option := arg.Option()

		if option.Name == name {
			return arg
		}
	}

	return nil
}

// Arg is an interface for defining arguments.
type Arg interface {
	flag.Value
	ArgOptionGetter
}

// ArgOption represents all possible underlying argument type.
type ArgOption struct {
	// typ represents the underlying argument type.
	typ types.Type

	// Name is the name of the argument.
	//
	// This is only useful in showing help or examples.
	Name string

	// Bind is a variable or field to set.
	Bind interface{}

	// Desc is the description of this argument.
	Desc string

	// Position is the index in the slice of arguments (e.g. 0).
	Position int

	// Separator is the separator to use when providing multiple arguments
	// for the same variable or field (i.e. a slice).
	Separator string

	// Required indicates whether this argument is required.
	Required bool
}

// Type returns the type of the arg.
func (o ArgOption) Type() types.Type {
	return o.typ
}

// SortArgOptionsByName sorts args by name.
type SortArgOptionsByName []ArgOption

func (n SortArgOptionsByName) Len() int           { return len(n) }
func (n SortArgOptionsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortArgOptionsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

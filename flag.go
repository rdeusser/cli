package cli

import (
	"fmt"
	"strings"

	"github.com/rdeusser/cli/internal/join"
)

var HelpFlag = &Flag[bool]{
	Name:      "help",
	Shorthand: "h",
	Desc:      "Print help information",
}

// Flags is a slice of flags represented as Options.
type Flags []option

func (flags Flags) Lookup(name string) option {
	for _, flag := range flags {
		if matchesFlag(name, flag) {
			return flag
		}
	}

	return nil
}

var _ option = (*Flag[bool])(nil)

// Flag is a generic type for defining flags with types constrained by Value.
type Flag[T Value] struct {
	Name      string
	Shorthand string
	Desc      string
	Separator byte   // only applies if the value is actually many
	Layout    string // only applies to time.Time values
	Default   T
	Value     *T
	EnvVar    EnvVar[T]
	Required  bool

	isSlice    bool
	hasBeenSet bool
}

// Init initializes the value of a flag.
func (f *Flag[T]) Init() error {
	if f.Value == nil {
		f.Value = new(T)
	}

	if len(f.Shorthand) > 1 {
		return ErrInvalidShorthand
	}

	if !isZeroValue(f.Default) {
		result, err := parseValue[T](f.Default, f.Separator, f.Layout)
		if err != nil {
			return err
		}

		*f.Value = result
	}

	if f.EnvVar.Name != "" {
		result, err := parseValue[T](f.EnvVar.Name, f.Separator, f.Layout)
		if err != nil {
			return err
		}

		*f.Value = result
	}

	return nil
}

// Set parses the value of s and sets the value according to the flags type.
func (f *Flag[T]) Set(s string) error {
	value, err := parseValue[T](s, f.Separator, f.Layout)
	if err != nil {
		return err
	}

	*f.Value = value
	f.isSlice = strings.HasPrefix(fmt.Sprint(value), "[")
	f.hasBeenSet = true

	return nil
}

// String returns the string form of the flags value.
func (f *Flag[T]) String() string {
	if f == nil || f.Value == nil {
		return ""
	}

	// If the flag type is a slice, we have to remove the brackets that
	// fmt.Sprint will add, and rejoin the string with the flags separator.
	return join.WithSeparator(trimBrackets(*f.Value), f.Separator)
}

// Options returns the common Options available to both flags and arguments.
func (f *Flag[T]) Options() Options {
	t := *new(T)
	if isZeroValue(f.Default) {
		f.Default = t
	}

	return Options{
		IsSlice:    f.isSlice,
		Name:       f.Name,
		Shorthand:  f.Shorthand,
		Desc:       f.Desc,
		Separator:  f.Separator,
		Layout:     f.Layout,
		Default:    f.Default,
		Value:      f.Value,
		EnvVar:     f.EnvVar,
		Required:   f.Required,
		HasBeenSet: f.hasBeenSet,
	}
}

// SortFlagsByName sorts flags by name.
type SortFlagsByName Flags

func (n SortFlagsByName) Len() int           { return len(n) }
func (n SortFlagsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortFlagsByName) Less(i, j int) bool { return n[i].Options().Name < n[j].Options().Name }

package cli

import (
	"fmt"
	"strings"
)

// Args is a slice of args.
type Args []option

// Lookup looks up an argument by it's position.
func (args Args) Lookup(position int) option {
	for i, arg := range args {
		if i == position {
			return arg
		}
	}

	return nil
}

var _ option = (*Arg[bool])(nil)

// Arg is a generic type for defining arguments with types constrained by Value.
type Arg[T Value] struct {
	Name     string
	Desc     string
	Layout   string // only applies to time.Time values
	Value    *T
	Required bool

	hasBeenSet bool
}

// Init initializes the value of an argument.
func (a *Arg[T]) Init() error {
	if a.Value == nil {
		a.Value = new(T)
	}

	return nil
}

// Set parses the value of s and sets the value according to the arguments type.
func (a *Arg[T]) Set(s string) error {
	value, err := parseValue[T](s, ' ', a.Layout)
	if err != nil {
		return err
	}

	*a.Value = value
	a.hasBeenSet = true

	return nil
}

// String returns the string form of the arguments value.
func (a *Arg[T]) String() string {
	if a == nil || a.Value == nil {
		return ""
	}

	return fmt.Sprint(*a.Value)
}

// Options returns the common Options available to both flags and arguments.
func (a *Arg[T]) Options() Options {
	t := *new(T)

	isSlice := false
	if strings.HasPrefix(a.String(), "[") {
		isSlice = true
	}

	return Options{
		IsSlice:    isSlice,
		Name:       a.Name,
		Desc:       a.Desc,
		Layout:     a.Layout,
		Default:    t,
		Value:      a.Value,
		Required:   a.Required,
		HasBeenSet: a.hasBeenSet,
	}
}

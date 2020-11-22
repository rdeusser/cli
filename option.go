package cli

import (
	"flag"
)

type TypeGetter interface {
	Type() OptionType
}

type OptionGetter interface {
	Option() (Option, error)
}

//go:generate stringer -type OptionType -linecomment
type OptionType int

const (
	Invalid  OptionType = iota // invalid
	Bool                       // bool
	String                     // string
	Int                        // int
	Float64                    // float64
	Duration                   // time.Duration
	Strings                    // []string
)

// Option represents a flag or argument.
type Option struct {
	optType OptionType

	Name      string
	Shorthand string
	Desc      string
	Default   string
	Value     flag.Value
	EnvVar    string
	Required  bool
}

func (o Option) Type() OptionType {
	return o.optType
}

func (o Option) HasBeenSet() bool {
	return o.Value.String() != ""
}

type SortOptionsByName []Option

func (n SortOptionsByName) Len() int           { return len(n) }
func (n SortOptionsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortOptionsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

package cli

import (
	"flag"
	"time"
)

type TypeGetter interface {
	GetType() OptionType
}

type OptionGetter interface {
	GetOption() (Option, error)
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
	Name      string
	Shorthand string
	Desc      string
	EnvVar    string
	Value     flag.Value
	Default   string
	Required  bool

	optType    OptionType
	hasBeenSet bool
}

func (o Option) GetType() OptionType {
	return o.optType
}

func (o Option) GetOption() (Option, error) {
	return o, nil
}

// BoolOpt represents a bool flag or argument.
type BoolOpt interface {
	TypeGetter
	OptionGetter

	value(*bool) (flag.Value, error)
}

// StringOpt represents a string flag or argument.
type StringOpt interface {
	TypeGetter
	OptionGetter

	value(*string) (flag.Value, error)
}

// IntOpt represents a int flag or argument.
type IntOpt interface {
	TypeGetter
	OptionGetter

	value(*int) (flag.Value, error)
}

// Float64Opt represents a float64 flag or argument.
type Float64Opt interface {
	TypeGetter
	OptionGetter

	value(*float64) (flag.Value, error)
}

// DurationOpt represents a time.Duration flag or argument.
type DurationOpt interface {
	TypeGetter
	OptionGetter

	value(*time.Duration) (flag.Value, error)
}

// StringsOpt represents a string slice flag or argument.
type StringsOpt interface {
	TypeGetter
	OptionGetter

	value(*[]string) (flag.Value, error)
}

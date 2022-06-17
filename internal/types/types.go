package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Getter interface {
	Type() Type
}

//go:generate stringer -type=Type -linecomment
type Type int

const (
	InvalidType     Type = iota // invalid
	BoolType                    // bool
	StringType                  // string
	IntType                     // int
	Float64Type                 // float64
	DurationType                // time.Duration
	StringSliceType             // []string
)

type Bool bool

func NewBool(into *bool, v bool) *Bool {
	if into == nil {
		into = new(bool)
	}

	*into = v
	return (*Bool)(into)
}

func (v *Bool) Set(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*v = Bool(b)
	return nil
}

func (v *Bool) String() string {
	return strconv.FormatBool(bool(*v))
}

func (v *Bool) Type() Type {
	return BoolType
}

type String string

func NewString(into *string, v string) *String {
	if into == nil {
		into = new(string)
	}

	*into = v
	return (*String)(into)
}

func (v *String) Set(s string) error {
	*v = String(s)
	return nil
}

func (v *String) String() string {
	return string(*v)
}

func (v *String) Type() Type {
	return StringType
}

type Int int

func NewInt(into *int, v int) *Int {
	if into == nil {
		into = new(int)
	}

	*into = v
	return (*Int)(into)
}

func (v *Int) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*v = Int(int(i))
	return nil
}

func (v *Int) String() string {
	return strconv.Itoa(int(*v))
}

func (v *Int) Type() Type {
	return IntType
}

type Float64 float64

func NewFloat64(into *float64, v float64) *Float64 {
	if into == nil {
		into = new(float64)
	}

	*into = v
	return (*Float64)(into)
}

func (v *Float64) Set(s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*v = Float64(f)
	return nil
}

func (v *Float64) String() string {
	return strconv.FormatFloat(float64(*v), 'g', -1, 64)
}

func (v *Float64) Type() Type {
	return Float64Type
}

type Duration time.Duration

func NewDuration(into *time.Duration, v time.Duration) *Duration {
	if into == nil {
		into = new(time.Duration)
	}

	*into = v
	return (*Duration)(into)
}

func (v *Duration) Set(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*v = Duration(d)
	return nil
}

func (v *Duration) String() string {
	return time.Duration(*v).String()
}

func (v *Duration) Type() Type {
	return DurationType
}

type StringSlice []string

func NewStringSlice(into *[]string, v []string) *StringSlice {
	if into == nil {
		into = new([]string)
	}

	*into = v
	return (*StringSlice)(into)
}

func (v *StringSlice) Set(s string) error {
	*v = append(*v, s)
	return nil
}

func (v *StringSlice) String() string {
	sb := new(strings.Builder)

	sb.WriteString("[")
	for i, s := range *v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%#v", s))
	}
	sb.WriteString("]")

	return sb.String()
}

func (v *StringSlice) Type() Type {
	return StringSliceType
}

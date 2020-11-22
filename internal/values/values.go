package values

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type BoolValue bool

func NewBool(into *bool, v bool) *BoolValue {
	if into == nil {
		into = new(bool)
	}

	*into = v
	return (*BoolValue)(into)
}

func (v *BoolValue) Set(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*v = BoolValue(b)
	return nil
}

func (v *BoolValue) String() string {
	return strconv.FormatBool(bool(*v))
}

type StringValue string

func NewString(into *string, v string) *StringValue {
	if into == nil {
		into = new(string)
	}

	*into = v
	return (*StringValue)(into)
}

func (v *StringValue) Set(s string) error {
	*v = StringValue(s)
	return nil
}

func (v *StringValue) String() string {
	return string(*v)
}

type IntValue int

func NewInt(into *int, v int) *IntValue {
	*into = v
	return (*IntValue)(into)
}

func (v *IntValue) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*v = IntValue(int(i))
	return nil
}

func (v *IntValue) String() string {
	return strconv.Itoa(int(*v))
}

type Float64Value float64

func NewFloat64(into *float64, v float64) *Float64Value {
	*into = v
	return (*Float64Value)(into)
}

func (v *Float64Value) Set(s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*v = Float64Value(f)
	return nil
}

func (v *Float64Value) String() string {
	return strconv.FormatFloat(float64(*v), 'g', -1, 64)
}

type DurationValue time.Duration

func NewDuration(into *time.Duration, v time.Duration) *DurationValue {
	*into = v
	return (*DurationValue)(into)
}

func (v *DurationValue) Set(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*v = DurationValue(d)
	return nil
}

func (v *DurationValue) String() string {
	return time.Duration(*v).String()
}

type StringsValue []string

func NewStrings(into *[]string, v []string) *StringsValue {
	*into = v
	return (*StringsValue)(into)
}

func (v *StringsValue) Set(s string) error {
	*v = append(*v, s)
	return nil
}

func (v *StringsValue) String() string {
	sb := new(strings.Builder)

	sb.WriteString("[")
	for i, s := range *v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%#v", s))
	}

	return sb.String()
}

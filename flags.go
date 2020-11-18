package cli

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rdeusser/cli/internal/values"
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

type Flag interface {
	TypeGetter
	OptionGetter
}

type BoolFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   bool
	Value     bool
	EnvVar    string
	Required  bool

	hasBeenSet *bool
}

func (f BoolFlag) GetType() OptionType {
	return Bool
}

func (f BoolFlag) GetOption() (Option, error) {
	value, err := f.value(nil)
	if err != nil {
		return Option{}, err
	}

	return Option{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   strconv.FormatBool(f.Default),
		Required:  f.Required,

		optType:    Bool,
		hasBeenSet: f.hasBeenSet,
	}, nil
}

func (f BoolFlag) value(into *bool) (flag.Value, error) {
	if into == nil {
		into = new(bool)
	}

	if len(f.Shorthand) > 1 {
		return nil, ErrInvalidShorthand
	}

	if !f.Default {
		f.Value = f.Default
		f.hasBeenSet = boolPtr(true)
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %q as a bool value for flag %s", v, f.Name)
		}

		f.Value = b
		f.hasBeenSet = boolPtr(true)
	}

	return values.NewBool(into, f.Value), nil
}

type StringFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   string
	Value     string
	EnvVar    string
	Required  bool

	hasBeenSet bool
}

func (f StringFlag) GetType() OptionType {
	return String
}

func (f StringFlag) GetOption() (Option, error) {
	value, err := f.value(nil)
	if err != nil {
		return Option{}, err
	}

	return Option{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   f.Default,
		Required:  f.Required,

		optType: String,
	}, nil
}

func (f StringFlag) value(into *string) (flag.Value, error) {
	if into == nil {
		into = new(string)
	}

	if len(f.Shorthand) > 1 {
		return nil, ErrInvalidShorthand
	}

	if f.Default != "" {
		f.Value = f.Default
		f.hasBeenSet = true
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		f.Value = v
		f.hasBeenSet = true
	}

	return values.NewString(into, f.Value), nil
}

type IntFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   int
	Value     int
	EnvVar    string
	Required  bool

	hasBeenSet bool
}

func (f IntFlag) GetType() OptionType {
	return Int
}

func (f IntFlag) GetOption() (Option, error) {
	value, err := f.value(nil)
	if err != nil {
		return Option{}, err
	}

	return Option{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   strconv.FormatInt(int64(f.Default), 0),
		Required:  f.Required,

		optType: Int,
	}, nil
}

func (f IntFlag) value(into *int) (flag.Value, error) {
	if into == nil {
		into = new(int)
	}

	if len(f.Shorthand) > 1 {
		return nil, ErrInvalidShorthand
	}

	if f.Default != 0 {
		f.Value = f.Default
		f.hasBeenSet = true
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		i, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %q as an int value for flag %s", v, f.Name)
		}

		f.Value = int(i)
		f.hasBeenSet = true
	}

	return values.NewInt(into, f.Value), nil
}

type Float64Flag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   float64
	Value     float64
	EnvVar    string
	Required  bool

	hasBeenSet bool
}

func (f Float64Flag) GetType() OptionType {
	return Float64
}

func (f Float64Flag) GetOption() (Option, error) {
	value, err := f.value(nil)
	if err != nil {
		return Option{}, err
	}

	return Option{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   strconv.FormatFloat(f.Default, 0, 0, 64),
		Required:  f.Required,

		optType: Float64,
	}, nil
}

func (f Float64Flag) value(into *float64) (flag.Value, error) {
	if into == nil {
		into = new(float64)
	}

	if len(f.Shorthand) > 1 {
		return nil, ErrInvalidShorthand
	}

	if f.Default != 0.0 {
		f.Value = f.Default
		f.hasBeenSet = true
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		f64, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %q as a float64 value for flag %s", v, f.Name)
		}

		f.Value = f64
		f.hasBeenSet = true
	}

	return values.NewFloat64(into, f.Value), nil
}

type DurationFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   time.Duration
	Value     time.Duration
	EnvVar    string
	Required  bool

	hasBeenSet bool
}

func (f DurationFlag) GetType() OptionType {
	return Duration
}

func (f DurationFlag) GetOption() (Option, error) {
	value, err := f.value(nil)
	if err != nil {
		return Option{}, err
	}

	return Option{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   f.Default.String(),
		Required:  f.Required,

		optType: Duration,
	}, nil
}

func (f DurationFlag) value(into *time.Duration) (flag.Value, error) {
	if into == nil {
		into = new(time.Duration)
	}

	if len(f.Shorthand) > 1 {
		return nil, ErrInvalidShorthand
	}

	if f.Default != time.Duration(0) {
		f.Value = f.Default
		f.hasBeenSet = true
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %q as a time.Duration value for flag %s", v, f.Name)
		}

		f.Value = d
		f.hasBeenSet = true
	}

	return values.NewDuration(into, f.Value), nil
}

type StringsFlag struct {
	Name      string
	Shorthand string
	Desc      string
	Default   []string
	Value     []string
	EnvVar    string
	Required  bool

	hasBeenSet bool
}

func (f StringsFlag) GetType() OptionType {
	return Strings
}

func (f StringsFlag) GetOption() (Option, error) {
	value, err := f.value(nil)
	if err != nil {
		return Option{}, err
	}

	return Option{
		Name:      f.Name,
		Shorthand: f.Shorthand,
		Desc:      f.Desc,
		EnvVar:    f.EnvVar,
		Value:     value,
		Default:   strings.Join(f.Default, " "),
		Required:  f.Required,

		optType: Strings,
	}, nil
}

func (f StringsFlag) value(into *[]string) (flag.Value, error) {
	if into == nil {
		into = new([]string)
	}

	if len(f.Shorthand) > 1 {
		return nil, ErrInvalidShorthand
	}

	if f.Default != nil {
		f.Value = f.Default
		f.hasBeenSet = true
	}

	envVar := strings.TrimSpace(f.EnvVar)
	if v, ok := os.LookupEnv(envVar); ok {
		f.Value = strings.Split(v, " ")
		f.hasBeenSet = true
	}

	return values.NewStrings(into, f.Value), nil
}

func boolPtr(b bool) *bool {
	return &b
}

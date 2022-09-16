package cli

type OptionSetter interface {
	SetOptions(flags Flags) error
}

type option interface {
	Init() error
	Set(string) error
	String() string
	Options() Options
}

type Options struct {
	IsSlice    bool
	Name       string
	Shorthand  string
	Desc       string
	Separator  byte
	Layout     string
	Default    any
	Value      any
	EnvVar     any
	Required   bool
	HasBeenSet bool
}

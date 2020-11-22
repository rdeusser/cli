package cli

import (
	"flag"
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
	flag.Value
	TypeGetter
	OptionGetter
}

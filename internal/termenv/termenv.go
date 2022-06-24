package termenv

import (
	"fmt"

	"github.com/muesli/termenv"
)

var NoColor = false

const (
	ColorRed    = termenv.ANSIRed
	ColorYellow = termenv.ANSIYellow
	ColorGreen  = termenv.ANSIGreen
)

func Colorize(color termenv.Color, format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)

	if NoColor {
		return s
	}

	return termenv.String(s).Foreground(color).String()
}

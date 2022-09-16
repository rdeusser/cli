package termenv

import (
	"fmt"

	"github.com/muesli/termenv"
)

const (
	ColorRed         = termenv.ANSIRed
	ColorYellow      = termenv.ANSIYellow
	ColorGreen       = termenv.ANSIGreen
	ColorBrightWhite = termenv.ANSIBrightWhite
)

type Color termenv.Color

func Red(format string, a ...any) string {
	return colorize(ColorRed, format, a...)
}

func Yellow(format string, a ...any) string {
	return colorize(ColorYellow, format, a...)
}

func Green(format string, a ...any) string {
	return colorize(ColorGreen, format, a...)
}

func BrightWhite(format string, a ...any) string {
	return colorizeBold(ColorBrightWhite, format, a...)
}

func colorize(color Color, format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return termenv.String(s).Foreground(color).String()
}

func colorizeBold(color Color, format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return termenv.String(s).Foreground(color).Bold().String()
}

package cli

import (
	"fmt"

	"github.com/muesli/termenv"
	"go.uber.org/atomic"
)

var (
	NoColor = atomic.NewBool(false)

	projectName = atomic.NewString("")
)

const (
	ColorRed    = termenv.ANSIRed
	ColorYellow = termenv.ANSIYellow
	ColorGreen  = termenv.ANSIGreen
)

func Output(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(nil, format, args...))
}

func Info(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(ColorGreen, "[INFO][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(ColorYellow, "[WARN][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(ColorRed, "[ERROR][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(ColorRed, "[FATAL][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	s := fmt.Sprint(colorize(ColorRed, "[PANIC][%s]:", projectName), fmt.Sprintf(format, args...))
	panic(s)

}

func colorize(color termenv.Color, format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)

	if NoColor.Load() {
		return s
	}

	return termenv.String(s).Foreground(color).String()
}

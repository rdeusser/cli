package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/muesli/termenv"
	"go.uber.org/atomic"
)

var stdout io.Writer

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
	fmt.Fprintln(stdout, colorize(nil, format, args...))
}

func Info(format string, args ...interface{}) {
	fmt.Fprintln(stdout, colorize(ColorGreen, "[INFO][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	fmt.Fprintln(stdout, colorize(ColorYellow, "[WARN][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Fprintln(stdout, colorize(ColorRed, "[ERROR][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	fmt.Fprintln(stdout, colorize(ColorRed, "[FATAL][%s]:", projectName), fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	s := fmt.Sprint(colorize(ColorRed, "[PANIC][%s]:", projectName), fmt.Sprintf(format, args...))
	panic(s)

}

func colorize(color termenv.Color, format string, args ...interface{}) string {
	if stdout == nil {
		stdout = os.Stdout
	}

	s := fmt.Sprintf(format, args...)

	if NoColor.Load() {
		return s
	}

	return termenv.String(s).Foreground(color).String()
}

package cli

import (
	"fmt"

	"github.com/muesli/termenv"
)

//go:generate gen-enum -type=OutputLevel
type OutputLevel int

const (
	DebugLevel OutputLevel = iota // name=DEBUG
	InfoLevel                     // name=INFO
	WarnLevel                     // name=WARN
	ErrorLevel                    // name=ERROR
	FatalLevel                    // name=FATAL
	PanicLevel                    // name=PANIC
)

const (
	ColorRed    = termenv.ANSIRed
	ColorYellow = termenv.ANSIYellow
	ColorGreen  = termenv.ANSIGreen
)

var (
	ProjectName = ""
	NoColor     = false
	EnableDebug = false
)

func Print(format string, args ...interface{}) {
	fmt.Fprintln(Output, colorize(nil, format, args...))
}

func Debug(format string, args ...interface{}) {
	if EnableDebug {
		fmt.Fprintln(Output, colorize(ColorYellow, "[%s][%s]:", DebugLevel, ProjectName), fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	fmt.Fprintln(Output, colorize(ColorGreen, "[%s][%s]:", InfoLevel, ProjectName), fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	fmt.Fprintln(Output, colorize(ColorYellow, "[%s][%s]:", WarnLevel, ProjectName), fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Fprintln(Output, colorize(ColorRed, "[%s][%s]:", ErrorLevel, ProjectName), fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	fmt.Fprintln(Output, colorize(ColorRed, "[%s][%s]:", FatalLevel, ProjectName), fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	s := fmt.Sprint(colorize(ColorRed, "[%s][%s]:", PanicLevel, ProjectName), fmt.Sprintf(format, args...))
	panic(s)
}

func colorize(color termenv.Color, format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)

	if NoColor {
		return s
	}

	return termenv.String(s).Foreground(color).String()
}

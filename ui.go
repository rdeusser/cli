package cli

import (
	"fmt"

	"github.com/muesli/termenv"
	"go.uber.org/atomic"
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

var NoColor = atomic.NewBool(false)

var (
	projectName = ""
	outputLevel = InfoLevel
)

func ProjectName() string {
	return projectName
}

func SetProjectName(name string) {
	mu.Lock()
	defer mu.Unlock()
	projectName = name
}

func OuptutLevel() OutputLevel {
	return outputLevel
}

func SetOutputLevel(level OutputLevel) {
	mu.Lock()
	defer mu.Unlock()
	outputLevel = level
}

func Output(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(nil, format, args...))
}

func Debug(format string, args ...interface{}) {
	if outputLevel <= DebugLevel {
		fmt.Fprintln(output, colorize(ColorYellow, "[%s][%s]:", DebugLevel, projectName), fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if outputLevel <= InfoLevel {
		fmt.Fprintln(output, colorize(ColorGreen, "[%s][%s]:", InfoLevel, projectName), fmt.Sprintf(format, args...))
	}
}

func Warn(format string, args ...interface{}) {
	if outputLevel <= WarnLevel {
		fmt.Fprintln(output, colorize(ColorYellow, "[%s][%s]:", WarnLevel, projectName), fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if outputLevel <= ErrorLevel {
		fmt.Fprintln(output, colorize(ColorRed, "[%s][%s]:", ErrorLevel, projectName), fmt.Sprintf(format, args...))
	}
}

func Fatal(format string, args ...interface{}) {
	if outputLevel <= FatalLevel {
		fmt.Fprintln(output, colorize(ColorRed, "[%s][%s]:", FatalLevel, projectName), fmt.Sprintf(format, args...))
	}
}

func Panic(format string, args ...interface{}) {
	if outputLevel <= PanicLevel {
		s := fmt.Sprint(colorize(ColorRed, "[%s][%s]:", PanicLevel, projectName), fmt.Sprintf(format, args...))
		panic(s)
	}
}

func colorize(color termenv.Color, format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)

	if NoColor.Load() {
		return s
	}

	return termenv.String(s).Foreground(color).String()
}

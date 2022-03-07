package cli

import (
	"fmt"

	"github.com/muesli/termenv"
	"go.uber.org/atomic"
)

//go:generate gen-enum -type=level
type OutputLevel int32

const (
	DebugLevel OutputLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

const (
	ColorRed    = termenv.ANSIRed
	ColorYellow = termenv.ANSIYellow
	ColorGreen  = termenv.ANSIGreen
)

var NoColor = atomic.NewBool(false)

var (
	projectName = atomic.NewString("")
	outputLevel = atomic.NewInt32(int32(InfoLevel))
)

func ProjectName() string {
	return projectName.Load()
}

func SetProjectName(name string) {
	projectName.Store(name)
}

func OuptutLevel() string {
	return ""
}

func SetOutputLevel(level OutputLevel) {
	outputLevel.Store(int32(level))
}

func Output(format string, args ...interface{}) {
	fmt.Fprintln(output, colorize(nil, format, args...))
}

func Debug(format string, args ...interface{}) {
	if outputLevel.Load() >= int32(DebugLevel) {
		fmt.Fprintln(output, colorize(ColorYellow, "[DEBUG][%s]:", projectName), fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if outputLevel.Load() >= int32(InfoLevel) {
		fmt.Fprintln(output, colorize(ColorGreen, "[INFO][%s]:", projectName), fmt.Sprintf(format, args...))
	}
}

func Warn(format string, args ...interface{}) {
	if outputLevel.Load() >= int32(WarnLevel) {
		fmt.Fprintln(output, colorize(ColorYellow, "[WARN][%s]:", projectName), fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if outputLevel.Load() >= int32(ErrorLevel) {
		fmt.Fprintln(output, colorize(ColorRed, "[ERROR][%s]:", projectName), fmt.Sprintf(format, args...))
	}
}

func Fatal(format string, args ...interface{}) {
	if outputLevel.Load() >= int32(FatalLevel) {
		fmt.Fprintln(output, colorize(ColorRed, "[FATAL][%s]:", projectName), fmt.Sprintf(format, args...))
	}
}

func Panic(format string, args ...interface{}) {
	if outputLevel.Load() >= int32(PanicLevel) {
		s := fmt.Sprint(colorize(ColorRed, "[PANIC][%s]:", projectName), fmt.Sprintf(format, args...))
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

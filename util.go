package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// rpad adds padding to the right side of a string.
func rpad(s string, count int) string {
	if count < 0 {
		count = 0
	}
	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", count))
}

// computePadding computes the padding needed for displaying usage text.
func computePadding(maxLen int, s string) int {
	return maxLen - len(s) + 4
}

// findMaxLength sorts a map of commands by their length and returns the length of the longest command name.
func findMaxLength(commands []*Command) int {
	if len(commands) == 0 {
		return 0
	}

	list := make([]int, 0, len(commands))

	for _, cmd := range commands {
		list = append(list, len(cmd.Name))
	}

	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(list)-1; i++ {
			if list[i+1] > list[i] {
				list[i+1], list[i] = list[i], list[i+1]
				swapped = true
			}
		}
	}

	return list[0]
}

// Following clap's lead here.
func good(format string, args ...interface{}) string {
	if format == "" {
		return ""
	}
	return color.GreenString(format, args...)
}

func warning(format string, args ...interface{}) string {
	if format == "" {
		return ""
	}
	return color.YellowString(format, args...)
}

func bad(format string, args ...interface{}) string {
	if format == "" {
		return ""
	}
	return color.RedString(format, args...)
}

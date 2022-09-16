package join

import "strings"

func Args(args []string) string {
	return strings.Join(args, " ")
}

func WithSeparator(arg string, sep byte) string {
	parts := strings.Split(arg, " ")
	return strings.Join(parts, string(sep))
}

package cli

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// ValueOf looks up the name of a flag and returns the value that it was set
// to. It's main use should be in the SetOptions method.
func ValueOf[T Value](flags Flags, name string) T {
	option := flags.Lookup(name)
	if option != nil && option.Options().HasBeenSet {
		flag := option.(*Flag[T])
		return *flag.Value
	}

	return *new(T)
}

func visit(fn func(*Command) error, commands []*Command) error {
	for _, cmd := range commands {
		if err := fn(cmd); err != nil {
			return err
		}
	}

	return nil
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func trimDash(s string) string {
	return strings.TrimLeft(s, "-")
}

func matchesFlag(arg string, opt option) bool {
	o := opt.Options()
	flag := trimDash(arg)

	if o.Name != flag && o.Shorthand != flag {
		return false
	}

	return true
}

func columnToSpace(column int) string {
	if column == 0 {
		return ""
	}

	return strings.Repeat(" ", column)
}

func lastIndex(s, substr string) int {
	tmp := strings.ReplaceAll(s, "\t", strings.Repeat(" ", 8))
	return strings.LastIndex(tmp, substr)
}

func splitBytes(s []byte, sep byte) [][]byte {
	if sep == 0 {
		return [][]byte{s}
	}

	return bytes.Split(s, []byte{sep})
}

func splitString(s string, sep byte) []string {
	if sep == 0 {
		return []string{s}
	}

	return strings.Split(s, string(sep))
}

func trimBrackets(s any) string {
	return strings.TrimFunc(fmt.Sprint(s), func(r rune) bool {
		return r == '[' || r == ']'
	})
}

// formatDesc formats a short description provided by a command, flag, or
// argument into the proper format.
//
// There's way too much variation in the way folks do this. Let me be clear: The
// first letter in the description is capital and there is no period at the end.
func formatDesc(s string) string {
	var sb strings.Builder

	if len(s) == 0 {
		return s
	}

	sb.WriteRune(unicode.ToUpper(rune(s[0])))
	sb.WriteString(strings.TrimSuffix(s[1:], "."))

	return sb.String()
}

type lineBuilder struct {
	currentLine strings.Builder
	sb          strings.Builder
}

func (lb *lineBuilder) Flush() {
	lb.sb.WriteString(lb.currentLine.String())
}

func (lb *lineBuilder) Write(s string) (n int, err error) {
	return lb.currentLine.WriteString(s)
}

func (lb *lineBuilder) NewLine() {
	lb.Flush()
	lb.sb.WriteString("\n")
	lb.currentLine.Reset()
}

func (lb *lineBuilder) CurrentLine() string {
	return lb.currentLine.String()
}

func (lb *lineBuilder) String() string {
	return lb.sb.String()
}

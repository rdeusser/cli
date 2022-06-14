package cli

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs = template.FuncMap{
	"trim":       trimRightSpace,
	"titleColor": titleColor,
}

func titleColor(s string) string {
	return colorize(ColorYellow, s)
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// renderTemplate executes the given template text on data, writing the result to w.
func renderTemplate(w io.Writer, text string, data interface{}) error {
	t, err := template.New("usage").Funcs(templateFuncs).Parse(text)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "usage", data)
}

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

var UsageTemplate = `{{ .Desc }}

{{ "USAGE:" | titleColor }}
    {{ .FullName }} {{ if .HasAvailableFlags }}[FLAGS]{{ end }}{{ if .HasAvailableCommands}} [command]{{ end }} [ARGUMENTS...]

{{- if .HasAvailableCommands }}

{{ "SUBCOMMANDS:" | titleColor }}
{{ .CommandHelp | trim }}
{{- end }}

{{- if .HasAvailableFlags }}

{{ "FLAGS:" | titleColor }}
{{ .FlagHelp | trim }}
{{- end }}

Use "{{ .FullName }} {{ if .HasAvailableCommands }}[command] {{ end }}--help" for more information about a command.
`

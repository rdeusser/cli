package cli

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/rdeusser/cli/internal/termenv"
)

var templateFuncs = template.FuncMap{
	"joinArgs": joinArgs,
	"title":    title,
	"trim":     trimRightSpace,
}

func joinArgs(argOptions []ArgOption, sep string) string {
	args := make([]string, 0)

	for _, option := range argOptions {
		args = append(args, fmt.Sprintf("<%s>", option.Name))
	}

	return strings.Join(args, sep)
}

func title(s string) string {
	return termenv.Colorize(termenv.ColorYellow, s)
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

var UsageTemplate = `{{ .Description }}

{{ "USAGE:" | title }}
    {{ .Command }} {{ if .HasAvailableFlags }}[flags]{{ end }}{{ if .HasAvailableCommands}} [command]{{ end }}{{ if .HasAvailableArgs}} {{ .Args }}{{ end }}

{{- if .HasAvailableCommands }}

{{ "SUBCOMMANDS:" | title }}
{{ .CommandHelp | trim }}
{{- end }}

{{- if .HasAvailableArgs }}

{{ "ARGS:" | title }}
{{ .ArgHelp | trim }}

{{- end }}

{{- if .HasAvailableFlags }}

{{ "FLAGS:" | title }}
{{ .FlagHelp | trim }}

{{- end }}

Use "{{ .Command }} {{ if .HasAvailableCommands }}[command] {{ end }}--help" for more information about a command.
`

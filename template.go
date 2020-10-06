package cli

import (
	"io"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs = template.FuncMap{
	"trim":    trimRightSpace,
	"warning": warning,
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// renderTemplate executes the given template text on data, writing the result to w.
func renderTemplate(w io.Writer, text string, data interface{}) error {
	t := template.New("usage").Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.ExecuteTemplate(w, "usage", data)
}

var UsageTemplate = `{{ .Desc }}

{{ "USAGE:" | warning}}
    {{ .Name }} {{ if .HasAvailableFlags }}[FLAGS]{{ end }}{{ if .HasAvailableCommands}} command{{ end }} [ARGUMENTS...]

{{- if .HasAvailableCommands }}

{{ "SUBCOMMANDS:" | warning }}
{{ .CommandHelp | trim }}
{{- end }}

{{- if .HasAvailableFlags }}

{{ "FLAGS:" | warning }}
{{ .FlagHelp | trim }}
{{- end }}

Use "{{ .Name }} {{ if .HasAvailableCommands }}[command] {{ end }}--help" for more information about a command.
`

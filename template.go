package cli

import (
	"io"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs = template.FuncMap{
	"trim":       trimRightSpace,
	"titleColor": warning,
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

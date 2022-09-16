package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs = template.FuncMap{
	"title": title,
}

var convertTemplate = `package cli

{{ range .SignedTypes }}
func convertInt64To{{ . | title }}(i []int64) []{{ . }} {
	var values []{{ . }}
	for _, v := range i {
		values = append(values, {{ . }}(v))
	}
	return values
}
{{ end }}

{{ range .UnsignedTypes }}
func convertUint64To{{ . | title }}(i []uint64) []{{ . }} {
	var values []{{ . }}
	for _, v := range i {
		values = append(values, {{ . }}(v))
	}
	return values
}
{{ end }}

{{ range .FloatTypes }}
func convertFloat64To{{ . | title }}(f []float64) []{{ . }} {
	var values []{{ . }}
	for _, v := range f {
		values = append(values, {{ . }}(v))
	}
	return values
}
{{ end }}

{{ range .ComplexTypes }}
func convertComplex128To{{ . | title }}(c []complex128) []{{ . }} {
	var values []{{ . }}
	for _, v := range c {
		values = append(values, {{ . }}(v))
	}
	return values
}
{{ end }}
`

func title(s string) string {
	var sb strings.Builder

	sb.WriteRune(unicode.ToUpper(rune(s[0])))
	sb.WriteString(s[1:])

	return sb.String()
}

func main() {
	filename := "type_converter_gen.go"

	tmpl, err := template.New("").Funcs(templateFuncs).Parse(convertTemplate)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		SignedTypes   []string
		UnsignedTypes []string
		FloatTypes    []string
		ComplexTypes  []string
	}{
		[]string{"int", "int8", "int16", "int32"},
		[]string{"uint", "uint8", "uint16", "uint32", "uintptr"},
		[]string{"float32"},
		[]string{"complex64"},
	}

	buf := bytes.Buffer{}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filename, src, 0o755); err != nil {
		log.Fatal(err)
	}
}

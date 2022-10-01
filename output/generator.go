package output

import (
	"bytes"
	"io"
	"text/template"

	"github.com/divpro/puml/parser"
)

var t *template.Template

var tpl = `@startuml
{{ range . }}
entity {{ .Name }} {
{{- range .Fields }}
	{{ if not .IsPtr  }}*{{ end -}}
	{{- .Name -}}
	{{- if not .IsEmbed }} : {{ end }}
	{{- if not .IsEmbed }}
		{{- if .IsSlice }}[]{{ end -}}
		{{- .Type -}}
	{{- end -}}
{{- end }}
}
{{ end }} 
@enduml
`

func init() {
	t = template.Must(template.New("").Parse(tpl))
}

func Out(structs []parser.PStruct) string {
	var buf bytes.Buffer
	err := t.Execute(io.Writer(&buf), structs)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

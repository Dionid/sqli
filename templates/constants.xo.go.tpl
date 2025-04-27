{{ define "debug" -}}
// {{- printf "Template params:\n" -}}
// {{ printf "%#v\n\n" . -}}
{{ end -}}

{{ define "constants" -}}
package {{ pkg }}

{{ $t := .Data }}

type TablesSt struct {
    {{ range $t -}}
        {{ .GoName }} string `json:"{{ .SQLName }}" db:"{{ .SQLName }}"`
    {{ end -}}
}

var Tables = TablesSt{
    {{ range $t -}}
        {{ .GoName }}: "{{ .SQLName }}",
    {{ end -}}
}

// Named "T" for shortness
var T = Tables

{{- end }}
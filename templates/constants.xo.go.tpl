{{ define "constants" -}}
package {{ pkg }}

// {{- printf "Template Data:\n" -}}
// {{ printf "%#v\n\n" . -}}

{{- $t := .Data -}}

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
package main

const tpl = `// Code generated by protoc-gen-go-enum. DO NOT EDIT.
// source: {{ .File.GeneratedFilenamePrefix }}.proto
// protoc-gen-go-enum version: {{ .Version }}

package {{ .File.GoPackageName }}

import (
	"database/sql/driver"
	"entgo.io/ent/schema/field"
	go_enum "github.com/devtoolskits/protoc-gen-go-enum/pkg"
)

var (
{{ range .Enums }}
	_ field.EnumValues = (*{{ .GoIdent.GoName }})(nil)
	_ field.ValueScanner = (*{{ .GoIdent.GoName }})(nil)
{{ end }}
)

{{ range .Enums }}
// Values implements ent/schema/field.EnumValues
func (x {{ .GoIdent.GoName }}) Values() []string {
	return go_enum.EnumMembers(x)
}

// Value implements sql/driver.Valuer
func (x {{ .GoIdent.GoName }}) Value() (driver.Value, error) {
	return x.String(), nil
}

// Scan implements sql.Scanner
func (x *{{ .GoIdent.GoName }}) Scan(src any) error {
	var s string
	switch v := src.(type) {
	case nil:
		return nil
	case string:
		s = v
	case []uint8:
		s = string(v)
	}

	n, ok := {{ .GoIdent.GoName }}_value[s]
	if ok {
		*x = {{ .GoIdent.GoName }}(n)
	}
	return nil
}
{{ end }}
`

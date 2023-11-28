package generate

import (
	"bytes"
	"fmt"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

type generateFileTemplateModel struct {
	PackageName string
	Type        typeModel
}
type typeModel struct {
	Name        string
	Annotations []string
	Fields      []fieldModel
}
type fieldModel struct {
	Name        string
	Type        string
	Annotations []string
}

var generateFileTemplate = template.Must(
	template.New("file").
		Funcs(sprig.TxtFuncMap()).
		Parse(`// GENERATED FILE DO NOT EDIT
		package {{ .PackageName }}

		import (
			"fmt"
		)

		{{template "load-func" .Type }}
		`))

var generateTypeLoadFuncTemplate = template.Must(
	generateFileTemplate.New("load-func").
		Parse(`func (s *{{ .Name }}) Factor3Load() error {
	fmt.Printf("type name: {{ .Name }}\n")
	fmt.Printf("annotations: {{range $a := .Annotations }}{{ $a }},{{ end }}\n")
	fmt.Printf("fields:\n")
	{{ range $field := .Fields }}
		fmt.Printf("\t-name={{ $field.Name }}, type={{ $field.Type }}\n")
		fmt.Printf("\t-annotations: {{range $a := $field.Annotations }}{{ $a }},{{ end }}\n")
	{{ end }}
	return nil
}
`))

func GenerateFileForType(t Type) ([]byte, error) {
	model := generateFileTemplateModel{
		PackageName: t.pkgName,
		Type: typeModel{
			Name:        t.name,
			Annotations: t.annotations,
			Fields:      nil,
		},
	}

	var fields []fieldModel
	for _, f := range t.fields {
		fields = append(fields, fieldModel{
			Name:        f.name,
			Type:        f.typ.String(),
			Annotations: f.annotations,
		})
	}
	model.Type.Fields = fields

	var results bytes.Buffer
	if err := generateFileTemplate.Execute(&results, model); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	return results.Bytes(), nil
}

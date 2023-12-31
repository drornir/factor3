package generate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/mattn/go-shellwords"
)

type generateFilesTemplateModel struct {
	PackageName string
	Type        typeModel
}
type typeModel struct {
	Config      configModel
	Name        string
	Annotations []string
	Fields      []fieldModel
}
type fieldModel struct {
	Name        string
	Type        string
	Annotations []string
	Doc         string
}
type configModel struct {
	ConfigFileName string
	EnvPrefix      string
}

var generateUtilsFileTemplate = template.Must(
	template.New("utils-file").
		Funcs(sprig.TxtFuncMap()).
		Parse(`// GENERATED FILE DO NOT EDIT
	package {{ .PackageName }}

	import (
		"encoding/json"
	)

	type zz_factor3_Config struct {
		Filename string
		EnvPrefix string
	}

	type zz_factor3_JSONer[T any] struct {t *T}
	func (s *zz_factor3_JSONer[T]) UnmarshalJSON(b []byte) error {
		return json.Unmarshal(b, s.t)
	}
	func (s *zz_factor3_JSONer[T]) MarshalJSON() ([]byte, error) {
		return json.Marshal(s.t)
	}
	
	`),
)

var generateTypeFileTemplate = template.Must(
	template.New("type-file").
		Funcs(sprig.TxtFuncMap()).
		Funcs(template.FuncMap{"shellwords": shellwords.Parse}).
		Parse(`// GENERATED FILE DO NOT EDIT
		package {{ .PackageName }}

		import (
			"encoding/json"
			"fmt"
			"strconv"
			"strings"
			"os"
			
			"gopkg.in/yaml.v3"
			"github.com/spf13/pflag"
		)

		var _ = strconv.ParseInt // just in case strconv is not used

		{{template "load-func" .Type }}
		`))

var generateTypeLoadFuncTemplate = template.Must(
	generateTypeFileTemplate.New("load-func").
		Parse(`func (self *{{ .Name }}) Factor3Load(argv []string) error {
	conf := zz_factor3_Config{}
	{{ if ne .Config.ConfigFileName "" }}conf.Filename = "{{ .Config.ConfigFileName }}"{{ end }}
	{{ if ne .Config.EnvPrefix "" }}conf.EnvPrefix = "{{ .Config.EnvPrefix }}_"{{ end }}

	{{ template "json_dec" . }}
	{{ template "env_dec" . }}
	{{ template "flag_dec" . }}

	if err := loadConfigFile(conf.Filename); err != nil {
		return fmt.Errorf("loading config from file %q: %w", conf.Filename, err)
	}
	if err := loadEnv(conf.EnvPrefix); err != nil {
		return fmt.Errorf("loading config from env: %w", err)
	}
	if err := parseFlags(argv); err != nil {
		return fmt.Errorf("loading config from pflags: %w", err)
	}
	
	return nil
}
`))

var _ = template.Must(
	generateTypeLoadFuncTemplate.New("json_dec").
		Parse(fmt.Sprintf(`
	loadConfigFile := func(filename string) error {
		type jsonStruct struct {
			{{- range $field := .Fields -}}
				{{ $field.Name }} {{ $field.Type }} %[1]cjson:"{{ $field.Name | snakecase }}"%[1]c
			{{ end -}}
			}
		var jsoner json.Unmarshaler
		if x, ok := interface{}(self).(json.Unmarshaler); ok {
			jsoner = x
		} else {
			jsoner = &zz_factor3_JSONer[jsonStruct]{t: (*jsonStruct)(self)}
		}
		file, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("opening file: %%w", err)
		}

		fileExt := filename[strings.LastIndex(filename, ".")+1:]
		switch fileExt {
		case "yaml", "yml":
			intoMap := make(map[string]interface{})
			err = yaml.Unmarshal(file, intoMap)
			if err != nil {
				break
			}
			intoJSON, e := json.Marshal(intoMap)
			if e != nil {
				err = e
				break
			}
			err = json.Unmarshal(intoJSON, jsoner)
		case "json":
			err = json.Unmarshal(file, jsoner)
		default:
			return fmt.Errorf("unsupported file type %%q", fileExt)
		}
		if err != nil {
			return fmt.Errorf("unmarshaling: %%w", err)
		}

		return nil
	}
`,
			0x60,
		)))
var _ = template.Must(
	generateTypeLoadFuncTemplate.New("env_dec").
		Parse(`
	loadEnv := func (prefix string) error {
		var s string
		{{- range $f := .Fields }}
			s = os.Getenv(prefix+"{{ $f.Name | snakecase | upper }}")

			{{- if eq $f.Type "string" }} 
			if s != "" {
				self.{{ $f.Name }} = s
			}
			{{- else if list "int" "int8" "int16" "int32" "int64" | has $f.Type }}
			if s != "" {
				if n, err := strconv.ParseInt(s, 10, {{ $f.Type | trimSuffix "int" | default "32" }}); err != nil {
					return fmt.Errorf("parsing \"{{ $f.Name }}\" as \"{{ $f.Type }}\": %w", err)
				} else {
				self.{{ $f.Name }} = {{ $f.Type }}(n)
				}
			}
			{{- else if list "uint" "uint8" "uint16" "uint32" "uint64" | has $f.Type }}
			if s != "" {
				if n, err := strconv.ParseUint(s, 10, {{ $f.Type | trimSuffix "uint" | default "32"}}); err != nil {
					return fmt.Errorf("parsing \"{{ $f.Name }}\" as \"{{ $f.Type }}\": %w", err)
				} else {
				self.{{ $f.Name }} = {{ $f.Type }}(n)
				}
			}
			{{- else if list "float32" "float64" | has $f.Type }} 
			if s != "" {
				if n, err := strconv.ParseFloat(s, 10, {{ $f.Type | trimSuffix "float" | default "32" }}); err != nil {
					return fmt.Errorf("parsing \"{{ $f.Name }}\" as \"{{ $f.Type }}\": %w", err)
				} else {
				self.{{ $f.Name }} = {{ $f.Type }}(n)
				}
			}
			{{- else if eq "bool" $f.Type }} 
			if s != "" {
				if b, err := strconv.ParseBool(s) ; err != nil {
					return fmt.Errorf("parsing \"{{ $f.Name }}\" as \"{{ $f.Type }}\": %w", err)
				} else {
					self.{{ $f.Name }} = b
				}
			}
			{{- else }}
				// {{ $f.Name }}: {{ $f.Type }} is not a valid type
			{{ end -}}
		{{- end }}
		return nil
	}
		`),
)
var _ = template.Must(
	generateTypeLoadFuncTemplate.New("flag_dec").
		Parse(fmt.Sprintf(`
	parseFlags := func(argv []string) error {
		if len(argv) == 0 {
			return nil
		}
		fset := pflag.NewFlagSet("{{ .Name }}", pflag.ContinueOnError)
		{{ range $f := .Fields -}}
			{{- range $a := $f.Annotations -}}
				{{- if $a | hasPrefix "%[1]spflag" -}}
				{{- $args := shellwords (trimPrefix "%[1]spflag" $a) -}}
				fset.{{ $f.Type | title }}VarP(&self.{{ $f.Name }}, "{{ first $args }}" , "{{ first (rest $args) | default "" }}", self.{{ $f.Name }}, "{{ $f.Doc }}")
				{{ end -}}
			{{- end -}}
		{{- end }}
		if err := fset.Parse(argv); err != nil {
			return fmt.Errorf("parsing flags: %%w", err)
		}
		return nil
	}	
`,
			FACTOR3_ANNOTATION_PREFIX,
		)))

func GenerateFilesForType(t Type) (map[string]string, error) {
	model := generateFilesTemplateModel{
		PackageName: t.pkgName,
		Type: typeModel{
			Config:      configModel(t.config),
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

	results := make(map[string]string)
	var b bytes.Buffer
	if err := generateUtilsFileTemplate.Execute(&b, model); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}
	results["zz_factor3_utils.go"] = b.String()
	b.Reset()

	if err := generateTypeFileTemplate.Execute(&b, model); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}
	typeName := strings.ToLower(model.Type.Name)
	results["zz_factor3_"+typeName+".go"] = b.String()
	b.Reset()

	return results, nil
}

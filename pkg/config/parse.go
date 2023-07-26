package config

type Schema struct {
	Fields []FieldSchema

	Env EnvSchema
}

type FieldSchema struct {
	Key   KeySchema
	Value ValueSchema
}

type KeySchema struct {
	JSONPath []string
}

type ValueSchema struct {
	Type ValueType
}

func ParseString(configStruct string) (Config, error) {
	return Config{
		Schema: Schema{
			Fields: nil,

			Env: EnvSchema{
				Usage: EnvUsage{
					shell: map[string]string{
						"USERNAME": "<string>",
						"PASSWORD": "<string>",
					},
				},
			},
		},
	}, nil
}

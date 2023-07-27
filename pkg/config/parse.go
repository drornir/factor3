package config

import (
	"fmt"
	"strings"
)

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
	ast, err := ASTFromString(configStruct)
	if err != nil {
		return Config{}, fmt.Errorf("generating abstract syntax tree from input: %w", err)
	}

	fieldSchemas, err := FieldSchemasFromAST(ast)
	if err != nil {
		return Config{}, fmt.Errorf(": %w", err)
	}

	envSchema := EnvSchemaFromFieldSchema(fieldSchemas)

	return Config{
		Schema: Schema{
			Fields: fieldSchemas,

			Env: envSchema,
		},
	}, nil
}

func FieldSchemasFromAST(tree AST) ([]FieldSchema, error) {
	mock := []FieldSchema{
		{
			Key: KeySchema{
				JSONPath: []string{"Username"},
			},
			Value: ValueSchema{
				Type: ValueTypeString,
			},
		},
		{
			Key: KeySchema{
				JSONPath: []string{"Password"},
			},
			Value: ValueSchema{
				Type: ValueTypeString,
			},
		},
	}

	return mock, nil
}

func EnvSchemaFromFieldSchema(fields []FieldSchema) EnvSchema {
	usage := make(map[string]string, len(fields))

	for _, field := range fields {
		k := strings.Join(field.Key.JSONPath, "_")
		k = strings.ToUpper(k)

		v := field.Value.Type.String()

		usage[k] = v
	}

	es := EnvSchema{
		Usage: usage,
	}

	return es

}

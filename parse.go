package factor3

type Schema struct {
	Fields []FieldSchema
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

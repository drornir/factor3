package parse

import (
	"go/ast"
	"log"
)

type FieldSchema struct {
	Name string
	Type string
}

func Field(field *ast.Field) []FieldSchema {
	var schemas []FieldSchema

	ast.Walk(fieldVisitor{
		schemas: &schemas,
	}, field)

	return schemas
}

type fieldVisitor struct {
	schemas *[]FieldSchema
}

func (av fieldVisitor) Visit(node ast.Node) ast.Visitor {
	log.Printf("fieldVisitor %#v", node)
	if node == nil {
		return nil
	}

	return av
}

package parse

import (
	"go/ast"
	"log"
)

type StructSchema struct {
	Fields []FieldSchema
}

func Struct(node *ast.StructType) StructSchema {
	ast.Walk(StructVisitor{}, node)

	return StructSchema{}
}

type StructVisitor struct{}

func (av StructVisitor) Visit(node ast.Node) ast.Visitor {
	log.Printf("StructVisitor %#v", node)

	return av
}

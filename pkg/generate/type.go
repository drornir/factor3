package generate

import (
	"go/ast"
	"go/types"
)

type UnparsedType struct {
	pkgID       string
	name        string
	doc         string
	annotations []string
	object      types.Object
	astTypeSpec *ast.TypeSpec
}

type Type struct {
	config      Config
	pkgName     string
	name        string
	annotations []string
	doc         string
	fields      []Field
}

type Field struct {
	parent      *Type
	name        string
	typ         types.Type
	annotations []string
	doc         string
}

type Config struct {
	ConfigFileName string
	EnvPrefix      string
}

package config

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/drornir/factor3/pkg/parse"
)

type AST struct {
}

func ASTFromString(configStruct string) (AST, error) {
	fset := token.NewFileSet()
	parserMode := parser.Trace | parser.ParseComments | parser.SkipObjectResolution
	astRoot, err := parser.ParseFile(fset, "inline", configStruct, parserMode)
	if err != nil {
		return AST{}, fmt.Errorf("parsing struct: %w", err)
	}

	ast.Walk(astVisitor{}, astRoot)

	return AST{}, nil
}

type astVisitor struct {
	Fields []FieldSchema
}

func (av astVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	fieldsSlice := []FieldSchema{}
	fieldsChan := make(chan FieldSchema)
	go func() {
		for f := range fieldsChan {
			fieldsSlice = append(fieldsSlice, f)
		}
	}()

	switch node.(type) {
	case *ast.GenDecl:
		return astVisitorGenDecl{C: fieldsChan}
	default:
		return av
	}
}

type astVisitorGenDecl struct {
	C chan<- FieldSchema
}

func (av astVisitorGenDecl) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.TypeSpec:
		log.Printf("TypeSpec %#v", n)
		return visitorFunc(func(node ast.Node) ast.Visitor {
			switch n := node.(type) {
			case *ast.StructType:
				return astVisitorStructType{C: av.C, Parent: n}
			default:
				return nil
			}
		})
	default:
		return nil
	}
}

type astVisitorStructType struct {
	C      chan<- FieldSchema
	Parent *ast.StructType
}

func (av astVisitorStructType) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.FieldList:
		log.Printf("StructType/FieldList %#v", n)
		return visitorFunc(func(node ast.Node) ast.Visitor {
			switch n := node.(type) {
			case *ast.Field:
				// return astVisitorField{
				// 	C:      av.C,
				// 	Parent: n,
				// }
				fs := parse.Field(n)
				for _, f := range fs {
					av.C <- FieldSchema{
						Key: JSONPath{f.Name},
						Value: ValueSchema{
							Type: ValueTypeString,
						},
					}
				}

				return nil
			default:
				return nil
			}
		})
	default:
		log.Printf("StructType/default %#v", n)
		return nil
	}
}

type astVisitorField struct {
	C      chan<- FieldSchema
	Parent *ast.Field
}

func (av astVisitorField) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return nil
	}

	var names []string // might have multiple identifiers e.g `X, Y string`
	var typ string

	// return visitorFunc(func(node ast.Node) ast.Visitor {

	// })

	// n := node.(*ast.Field)
	typeVisitor := visitorFunc(func(node ast.Node) ast.Visitor {
		if node == nil {
			return nil
		}

		switch n := node.(type) {
		case *ast.Ident:
			log.Printf("Field/Type/Ident %#v", node)
			typ = n.Name + typ
			return nil
		case *ast.SelectorExpr:
			var vf func(node ast.Node) ast.Visitor
			vf = func(node ast.Node) ast.Visitor {
				if node == nil {
					return nil
				}

				switch n := node.(type) {
				case *ast.Ident:
					typ = "." + n.Name + typ
					return visitorFunc(vf)
				default:
					log.Printf("Field/Type/SelectorExpr/default %#v", node)
					return nil
				}
			}
			return visitorFunc(vf)
		default:
			log.Printf("Field/Type/default %#v", node)
			return nil
		}
	})

	ast.Walk(typeVisitor, node)
	log.Printf("$$$$$$$$$$$$$ typ=%#v", typ)

	// if n.Names == nil {
	// 	names = []string{typ}
	// }

	_ = names
	//

	return nil
}

type visitorFuncWrapper struct {
	f func(node ast.Node) ast.Visitor
}

func (w visitorFuncWrapper) Visit(node ast.Node) ast.Visitor {
	return w.f(node)
}

func visitorFunc(f func(node ast.Node) ast.Visitor) ast.Visitor {
	return visitorFuncWrapper{
		f: f,
	}
}

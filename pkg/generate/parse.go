package generate

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

const FACTOR3_ANNOTATION_PREFIX = "//factor3:"

func (a *App) parseType(ut *UnparsedType) (Type, error) {
	var t Type
	t.name = ut.name
	t.pkgName = ut.object.Pkg().Name()
	t.annotations = ut.annotations

	ti := a.pkgs[ut.pkgID].TypesInfo.Types[ut.astTypeSpec.Type]
	switch uti := ti.Type.Underlying().(type) {
	case *types.Struct:
		for i := 0; i < uti.NumFields(); i++ {
			var fieldAnnotations []string

			field := uti.Field(i)
			if !field.Exported() {
				continue
			}

			ast.Inspect(ut.astTypeSpec, func(node ast.Node) bool {
				if node == nil {
					return true
				}

				switch n := node.(type) {
				case *ast.Field:
					if n.Doc == nil {
						return true
					}

					var found bool
					for _, nameIdent := range n.Names {
						if nameIdent.Name == field.Name() {
							found = true
							break
						}
					}
					if !found {
						return true
					}
					for _, c := range n.Doc.List {
						trimmed := strings.TrimSpace(c.Text)
						if strings.HasPrefix(trimmed, FACTOR3_ANNOTATION_PREFIX) {
							fieldAnnotations = append(fieldAnnotations, trimmed)
						}
					}
					return false
				default:
					return true
				}
			})

			outField := Field{
				parent:      &t,
				name:        field.Name(),
				typ:         field.Type().Underlying(),
				annotations: fieldAnnotations,
			}

			t.fields = append(t.fields, outField)
		}
	default:
		return t, fmt.Errorf("parsing type: only struct types are supported for generation: got %q", uti.String())
	}

	var genAnnotation string
	prefix := FACTOR3_ANNOTATION_PREFIX + "generate "
	for _, a := range t.annotations {
		if strings.HasPrefix(a, prefix) {
			genAnnotation = strings.TrimPrefix(a, prefix)
			break
		}
	}

	c, _ := parseConfig(genAnnotation) // ignoring error since the empty config is a good handling
	t.config = c

	return t, nil
}

func (a *App) generateCode(t Type) (map[string]string, error) {
	return GenerateFilesForType(t)
}

func (a *App) getAnnotatedTypes(pkgID string) map[string]*UnparsedType {
	results := make(map[string]*UnparsedType)
	visitor := &annotationsVisitor{app: a, pkgID: pkgID, results: results}
	for _, f := range a.pkgs[pkgID].Syntax {
		ast.Walk(visitor, f)
	}

	// filter unannotated types
	for key, result := range results {
		var found bool
		for _, a := range result.annotations {
			if strings.HasPrefix(a, FACTOR3_ANNOTATION_PREFIX) {
				found = true
				break
			}
		}

		if !found {
			delete(results, key)
		}
	}

	return results
}

type annotationsVisitor struct {
	app     *App
	pkgID   string
	results map[string]*UnparsedType

	lastDoc []*ast.Comment
}

func (v *annotationsVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.FuncDecl:
		return nil
	case *ast.GenDecl:
		if n.Tok != token.TYPE {
			return nil
		}
		if n.Doc != nil {
			v.lastDoc = n.Doc.List
		} else {
			v.lastDoc = nil
		}
		return v

	case *ast.TypeSpec:
		fmt.Println("v.lastDoc", v.lastDoc)
		var annotations []string
		for _, line := range v.lastDoc {
			if strings.HasPrefix(line.Text, FACTOR3_ANNOTATION_PREFIX) {
				annotations = append(annotations, line.Text)
			}
		}
		fmt.Println("annotations", annotations)

		pkg := v.app.pkgs[v.pkgID]
		object := pkg.Types.Scope().Lookup(n.Name.String())
		if object == nil {
			panic(fmt.Errorf("can't find types.Object for %q", n.Name.String()))
		}

		v.results[n.Name.String()] = &UnparsedType{
			pkgID:       v.pkgID,
			name:        n.Name.String(),
			annotations: annotations,
			object:      object,
			astTypeSpec: n,
		}
		return nil
	default:
		return v
	}
}

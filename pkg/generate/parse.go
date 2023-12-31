package generate

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

const FACTOR3_ANNOTATION_PREFIX = "//factor3:"
const FACTOR3_ANNOTATION_PREFIX_GENERATE = FACTOR3_ANNOTATION_PREFIX + "generate "
const FACTOR3_ANNOTATION_PREFIX_PFLAG = FACTOR3_ANNOTATION_PREFIX + "pflag "

func (a *App) parseType(ut *UnparsedType) (Type, error) {
	var t Type
	t.name = ut.name
	t.pkgName = ut.object.Pkg().Name()
	t.annotations = ut.annotations
	t.doc = ut.doc

	ti := a.pkgs[ut.pkgID].TypesInfo.Types[ut.astTypeSpec.Type]
	switch uti := ti.Type.Underlying().(type) {
	case *types.Struct:
		for i := 0; i < uti.NumFields(); i++ {
			var (
				fieldAnnotations []string
				fieldDocs        string
			)

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
					if n.Doc != nil {
						for _, c := range n.Doc.List {
							trimmed := strings.TrimSpace(c.Text)
							if strings.HasPrefix(trimmed, FACTOR3_ANNOTATION_PREFIX) {
								fieldAnnotations = append(fieldAnnotations, trimmed)
							}
						}
						fieldDocs = n.Doc.Text()
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
				doc:         fieldDocs,
			}

			t.fields = append(t.fields, outField)
		}
	default:
		return t, fmt.Errorf("parsing type: only struct types are supported for generation: got %q", uti.String())
	}

	var genAnnotation string
	for _, a := range t.annotations {
		if strings.HasPrefix(a, FACTOR3_ANNOTATION_PREFIX_GENERATE) {
			genAnnotation = strings.TrimPrefix(a, FACTOR3_ANNOTATION_PREFIX_GENERATE)
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
			if strings.HasPrefix(a, FACTOR3_ANNOTATION_PREFIX_GENERATE) {
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

	lastDoc *ast.CommentGroup
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

		v.lastDoc = n.Doc
		return v

	case *ast.TypeSpec:
		var annotations []string
		if v.lastDoc != nil {
			for _, line := range v.lastDoc.List {
				if strings.HasPrefix(line.Text, FACTOR3_ANNOTATION_PREFIX) {
					annotations = append(annotations, line.Text)
				}
			}
		}

		pkg := v.app.pkgs[v.pkgID]
		object := pkg.Types.Scope().Lookup(n.Name.String())
		if object == nil {
			panic(fmt.Errorf("can't find types.Object for %q", n.Name.String()))
		}

		v.results[n.Name.String()] = &UnparsedType{
			pkgID:       v.pkgID,
			name:        n.Name.String(),
			annotations: annotations,
			doc:         v.lastDoc.Text(),
			object:      object,
			astTypeSpec: n,
		}
		return nil
	default:
		return v
	}
}

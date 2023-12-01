package generate

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/packages"
)

const FACTOR3_ANNOTATION_PREFIX = "//factor3:"

type App struct {
	workDir string
	pkgs    map[string]*packages.Package
}

func New(workDir string) (*App, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("initializing new app: %w", err)
	}

	pkgs, err := LoadPackages(workDir)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &App{
		workDir: workDir,
		pkgs:    pkgs,
	}, nil
}

func (a *App) Generate() error {
	for pkgID := range a.pkgs {
		if err := a.generatePackage(pkgID); err != nil {
			return fmt.Errorf("generating code (package %q): %w", pkgID, err)
		}
	}
	return nil
}

func (a *App) generatePackage(pkgID string) error {
	unparsedTypes := a.getAnnotatedTypes(pkgID)

	allCode := make(map[string]string)
	for _, ut := range unparsedTypes {
		t, err := a.parseType(ut)
		if err != nil {
			return fmt.Errorf("generating in package %q: %w", pkgID, err)
		}

		code, err := a.generateCode(t)
		if err != nil {
			return fmt.Errorf("generating code for type %q: %w", t.name, err)
		}
		for k, v := range code {
			allCode[k] = v
		}
	}

	for filename, content := range allCode {
		fmt.Printf("Generated Code for %q (pkg=%q)\n", filename, pkgID)
		fmt.Println(string(content))

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing file %q in package %q: %w", filename, pkgID, err)
		}
	}

	if err := exec.Command("go", "fmt", a.workDir).Run(); err != nil {
		return fmt.Errorf("'go fmt' error: %w", err)
	}
	if err := exec.Command("go", "get", a.workDir).Run(); err != nil {
		return fmt.Errorf("'go get' error: %w", err)
	}

	return nil
}

func (a *App) parseType(ut *UnparsedType) (Type, error) {
	var t Type
	t.name = ut.name
	t.pkgName = ut.object.Pkg().Name()
	t.annotations = ut.annotations

	ti := a.pkgs[ut.pkgID].TypesInfo.Types[ut.astTypeSpec.Type]
	switch uti := ti.Type.Underlying().(type) {
	case *types.Struct:
		var fieldAnnotations []string
		for i := 0; i < uti.NumFields(); i++ {
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

	// fmt.Println("##########filename", c.ConfigFileName)

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
		var annotations []string
		for _, line := range v.lastDoc {
			if strings.HasPrefix(line.Text, FACTOR3_ANNOTATION_PREFIX) {
				annotations = append(annotations, line.Text)
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
			object:      object,
			astTypeSpec: n,
		}
		return nil
	default:
		return v
	}
}

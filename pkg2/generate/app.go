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
			return fmt.Errorf("generating code: %w", err)
		}
		fmt.Println(a.workDir)
		if err := exec.Command("go", "fmt", a.workDir).Run(); err != nil {
			return fmt.Errorf("format: %w", err)
		}
	}
	return nil
}

func (a *App) generatePackage(pkgID string) error {
	unparsedTypes := a.getAnnotatedTypes(pkgID)

	for _, ut := range unparsedTypes {
		t, err := a.parseType(ut)
		if err != nil {
			return fmt.Errorf("generating in package %q: %w", pkgID, err)
		}

		code, err := a.generateCode(t)
		if err != nil {
			return fmt.Errorf("generating in package %q: %w", pkgID, err)
		}

		// TODO write to file
		fmt.Printf("Generated Code for %q (pkg=%q)\n", t.name, pkgID)
		fmt.Println(string(code))

		filename := fmt.Sprintf("zz_factor3_%s.go", strings.ToLower(t.name))
		if err := os.WriteFile(filename, code, 0644); err != nil {
			return fmt.Errorf("writing file %q in package %q: %w", filename, pkgID, err)
		}
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

	return t, nil
}

func (a *App) generateCode(t Type) ([]byte, error) {
	return GenerateFileForType(t)
}

func (a *App) getAnnotatedTypes(pkgID string) map[string]*UnparsedType {
	results := make(map[string]*UnparsedType)
	visitor := &annotationsVisitor{app: a, pkgID: pkgID, results: results}
	for _, f := range a.pkgs[pkgID].Syntax {
		ast.Walk(visitor, f)
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
	case *ast.GenDecl:
		if n.Tok != token.TYPE {
			return nil
		}
		v.lastDoc = n.Doc.List
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
	default:
		return v
	}

	return nil
}

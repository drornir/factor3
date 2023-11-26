package generate

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
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

	log.Printf("loaded %d pkgs", len(pkgs))

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

	log.Printf("unparsed types #%d", len(unparsedTypes))

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
		fmt.Println(code)

		filename := fmt.Sprintf("zz_factor3_%s.go", strings.ToLower(t.name))
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
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
		for i := 0; i < uti.NumFields(); i++ {
			field := uti.Field(i)
			if !field.Exported() {
				continue
			}

			outField := Field{
				parent: &t,
				name:   field.Name(),
				typ:    field.Type().Underlying(),
			}

			t.fields = append(t.fields, outField)
		}
	default:
		return t, fmt.Errorf("parsing type: only struct types are supported for generation: got %q", uti.String())
	}

	return t, nil
}

func (a *App) generateCode(t Type) (string, error) {
	fieldsCode := ""
	for _, f := range t.fields {
		fieldsCode += fmt.Sprintf(`fmt.Printf("\t-%s (%s)\n")`+"\n", f.name, f.typ.String())
	}

	funcCode := fmt.Sprintf(`func (s *%[1]s) Factor3Load() error {
		fmt.Printf("type name: %%s\n", "%[1]s")
		fmt.Printf("annotations: %[2]s\n")
		fmt.Printf("fields:\n")
		%[3]s
		return nil
	}
	`,
		t.name, t.annotations, fieldsCode)

	final := fmt.Sprintf(`package %[1]s

	import (
		"fmt"
	)

	%[2]s
	`, t.pkgName, funcCode,
	)

	return final, nil
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
	log.Printf("visit: %T", node)

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
			log.Printf("doc line: %q", line.Text)
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

package generate

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/tools/go/packages"
)

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
		// fmt.Printf("Generated Code for %q (pkg=%q)\n", filename, pkgID)
		// fmt.Println(string(content))

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

package generate

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedExportFile |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes |
	packages.NeedModule |
	packages.NeedEmbedFiles |
	packages.NeedEmbedPatterns

func LoadPackages(patterns ...string) (map[string]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:    loadMode,
		Context: context.Background(),
		Logf: func(format string, args ...any) {
			slog.Debug(fmt.Sprintf(format, args...))
		},
		BuildFlags: []string{},
		Tests:      false,
	}

	pkgSlice, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("loading package(s) %q: %w", strings.Join(patterns, ", "), err)
	}

	pkgs := make(map[string]*packages.Package)
	for _, pkg := range pkgSlice {
		pkgs[pkg.ID] = pkg
	}

	return pkgs, nil
}

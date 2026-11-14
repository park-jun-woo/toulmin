//ff:func feature=tangl type=util control=sequence
//ff:what TestLoadRegistryErrors — tests loadAST-error, removed-cwd Abs failure, and alias-resolution error propagation
package registry

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestLoadRegistryErrors(t *testing.T) {
	t.Run("loadAST error propagates", func(t *testing.T) {
		_, err := loadRegistry("/base/missing.geul", fakeLoader(map[string]*ast.Document{}, map[string]int{}))
		if err == nil || !strings.Contains(err.Error(), "no document") {
			t.Fatalf("expected a loader error, got %v", err)
		}
	})

	t.Run("relative root with removed cwd fails Abs", func(t *testing.T) {
		withRemovedCwd(t)
		if _, err := loadRegistry("rel/a.geul", fakeLoader(map[string]*ast.Document{}, map[string]int{})); err == nil {
			t.Fatal("expected filepath.Abs to fail on a removed working directory")
		}
	})

	t.Run("alias resolution error propagates", func(t *testing.T) {
		// The loader stamps a *relative* document Path, so the alias's
		// resolveCheckingAliasPath is forced through filepath.Abs; with the
		// working directory removed that Abs fails and the error propagates.
		loader := func(string) (*ast.Document, error) {
			return &ast.Document{
				Path:  "rel/a.geul",
				Sees:  []ast.See{{Alias: "b", Package: "b"}},
				Cases: []ast.Case{aliasCheckNode("ca", "b", "cb")},
			}, nil
		}
		withRemovedCwd(t)
		if _, err := loadRegistry("/abs/a.geul", loader); err == nil {
			t.Fatal("expected the alias resolution error to propagate")
		}
	})
}

//ff:func feature=tangl type=codegen control=sequence
//ff:what TestResolveAliasImport — covers checking and using alias resolution, skips, and error paths
package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestResolveAliasImport(t *testing.T) {
	t.Run("checking alias resolves to its sibling import path", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/m\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		doc := &ast.Document{Path: filepath.Join(root, "a.geul"), Sees: []ast.See{{Alias: "b", Package: "b"}}}
		spec, ok, err := resolveAliasImport(doc, "b", true)
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if spec.Alias != "b" || spec.Path != "example.com/m" {
			t.Errorf("spec = %+v, want {b example.com/m}", spec)
		}
	})

	t.Run("checking alias without a See is skipped", func(t *testing.T) {
		doc := &ast.Document{Path: "/x/a.geul"}
		if _, ok, err := resolveAliasImport(doc, "missing", true); err != nil || ok {
			t.Fatalf("expected ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("checking alias resolve error propagates", func(t *testing.T) {
		doc := &ast.Document{Path: "rel/a.geul", Sees: []ast.See{{Alias: "b", Package: "b"}}}
		withRemovedCwd(t)
		if _, _, err := resolveAliasImport(doc, "b", true); err == nil {
			t.Fatal("expected the resolve error to propagate")
		}
	})

	t.Run("checking alias derive error propagates", func(t *testing.T) {
		// The sibling path resolves fine but no go.mod encloses it, so
		// deriveGoImportPath fails.
		doc := &ast.Document{Path: "/nonexistent-module-root-xyz/a.geul", Sees: []ast.See{{Alias: "b", Package: "b"}}}
		if _, _, err := resolveAliasImport(doc, "b", true); err == nil {
			t.Fatal("expected the derive error to propagate")
		}
	})

	t.Run("using alias with a See uses its package path", func(t *testing.T) {
		doc := &ast.Document{Sees: []ast.See{{Alias: "log", Package: "tangl/log"}}}
		spec, ok, err := resolveAliasImport(doc, "log", false)
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if want := "github.com/park-jun-woo/toulmin/go/pkg/tangl/log"; spec.Path != want {
			t.Errorf("path = %q, want %q", spec.Path, want)
		}
	})

	t.Run("using alias without a See falls back to tangl/<alias>", func(t *testing.T) {
		spec, ok, err := resolveAliasImport(&ast.Document{}, "foo", false)
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if want := "github.com/park-jun-woo/toulmin/go/pkg/tangl/foo"; spec.Path != want {
			t.Errorf("path = %q, want %q", spec.Path, want)
		}
	})
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestFindSee — covers matched and unmatched See-alias lookup
package registry

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestFindSee(t *testing.T) {
	d := &ast.Document{Sees: []ast.See{
		{Alias: "a", Package: "pkg/a"},
		{Alias: "b", Package: "pkg/b"},
	}}

	t.Run("alias found", func(t *testing.T) {
		see, ok := findSee(d, "b")
		if !ok || see.Package != "pkg/b" {
			t.Fatalf("expected pkg/b, got %q ok=%v", see.Package, ok)
		}
	})

	t.Run("alias not found", func(t *testing.T) {
		if _, ok := findSee(d, "missing"); ok {
			t.Fatal("expected ok=false for an undeclared alias")
		}
	})
}

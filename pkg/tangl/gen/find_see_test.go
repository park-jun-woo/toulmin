//ff:func feature=tangl type=codegen control=sequence
//ff:what TestFindSee — tests findSee for found and not-found branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestFindSee(t *testing.T) {
	doc := &ast.Document{
		Sees: []ast.See{
			{Alias: "a", Package: "pkg/a"},
			{Alias: "b", Package: "pkg/b"},
		},
	}

	see, ok := findSee(doc, "b")
	if !ok || see.Package != "pkg/b" {
		t.Errorf("expected to find alias %q, got %+v, %v", "b", see, ok)
	}

	if _, ok := findSee(doc, "missing"); ok {
		t.Error("expected not found for missing alias")
	}

	if _, ok := findSee(&ast.Document{}, "anything"); ok {
		t.Error("expected not found for empty Sees")
	}
}

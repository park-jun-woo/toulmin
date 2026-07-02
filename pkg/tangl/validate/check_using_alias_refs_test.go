//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckUsingAliasRefs — tests checkUsingAliasRefs for no-cases, nil-using, empty-alias, declared-alias, and undeclared-alias branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckUsingAliasRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkUsingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("NilUsing", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Using: nil},
				}},
			},
		}
		errs := checkUsingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("EmptyAlias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Using: &ast.Ref{Alias: "", Name: "fn"}},
				}},
			},
		}
		errs := checkUsingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "pkg", Package: "some/pkg"},
			},
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Using: &ast.Ref{Alias: "pkg", Name: "fn"}},
				}},
			},
		}
		errs := checkUsingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("UndeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Using: &ast.Ref{Alias: "missing", Name: "fn"}, Line: 7},
				}},
			},
		}
		errs := checkUsingAliasRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "undeclared package alias") {
			t.Errorf("expected undeclared package alias error, got %v", errs[0])
		}
	})
}

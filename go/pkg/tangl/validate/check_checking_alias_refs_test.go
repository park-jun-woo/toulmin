//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckCheckingAliasRefs — tests checkCheckingAliasRefs for no-cases, nil-checking, empty-alias, declared-alias, and undeclared-alias branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckCheckingAliasRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkCheckingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("NilChecking", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Checking: nil},
				}},
			},
		}
		errs := checkCheckingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("EmptyAlias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Checking: &ast.Ref{Alias: "", Name: "y"}},
				}},
			},
		}
		errs := checkCheckingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "gate", Package: "some/pkg"},
			},
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Checking: &ast.Ref{Alias: "gate", Name: "can pass"}},
				}},
			},
		}
		errs := checkCheckingAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("UndeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", Checking: &ast.Ref{Alias: "missing", Name: "can pass"}, Line: 7},
				}},
			},
		}
		errs := checkCheckingAliasRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "undeclared package alias") {
			t.Errorf("expected undeclared package alias error, got %v", errs[0])
		}
	})
}

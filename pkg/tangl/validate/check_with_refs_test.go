//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckWithRefs — tests checkWithRefs for no-cases, defined-term, and undefined-term branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckWithRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkWithRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DefinedTerm", func(t *testing.T) {
		doc := &ast.Document{
			Defs: []ast.Definition{
				{Name: "term1"},
			},
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", With: []string{"term1"}},
				}},
			},
		}
		errs := checkWithRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("UndefinedTerm", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Nodes: []ast.Node{
					{Name: "n1", With: []string{"missing"}, Line: 9},
				}},
			},
		}
		errs := checkWithRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "references undefined term") {
			t.Errorf("expected undefined term error, got %v", errs[0])
		}
	})
}

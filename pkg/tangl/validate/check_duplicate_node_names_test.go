//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicateNodeNames — tests checkDuplicateNodeNames for empty, unique, and duplicate-name branches across cases via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckDuplicateNodeNames(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkDuplicateNodeNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Line: 1},
						{Name: "b", Line: 2},
					},
				},
			},
		}
		errs := checkDuplicateNodeNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DuplicateAcrossMultipleCases", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Line: 1},
						{Name: "a", Line: 3},
					},
				},
				{
					Name: "y",
					Nodes: []ast.Node{
						{Name: "b", Line: 5},
					},
				},
			},
		}
		errs := checkDuplicateNodeNames(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), `node name in case "x"`) {
			t.Errorf("expected node name in case error, got %v", errs[0])
		}
	})
}

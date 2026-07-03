//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckSupersedeRefs — tests checkSupersedeRefs for valid refs, missing superseder, missing target, and both-missing branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckSupersedeRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkSupersedeRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "a"}, {Name: "b"}},
					Supersedes: []ast.Supersede{
						{Superseder: "a", Superseded: "b", Line: 1},
					},
				},
			},
		}
		errs := checkSupersedeRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("MissingTarget", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "a"}},
					Supersedes: []ast.Supersede{
						{Superseder: "a", Superseded: "missing", Line: 2},
					},
				},
			},
		}
		errs := checkSupersedeRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "target") {
			t.Errorf("expected target error, got %v", errs[0])
		}
	})

	t.Run("MissingSuperseder", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "b"}},
					Supersedes: []ast.Supersede{
						{Superseder: "missing", Superseded: "b", Line: 3},
					},
				},
			},
		}
		errs := checkSupersedeRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "superseder") {
			t.Errorf("expected superseder error, got %v", errs[0])
		}
	})

	t.Run("BothMissing", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{},
					Supersedes: []ast.Supersede{
						{Superseder: "missing-a", Superseded: "missing-b", Line: 4},
					},
				},
			},
		}
		errs := checkSupersedeRefs(doc)
		if len(errs) != 2 {
			t.Fatalf("expected 2 errors, got %v", errs)
		}
	})
}

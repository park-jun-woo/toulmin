//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckAttackRefs — tests checkAttackRefs for valid refs, missing target, missing attacker, and both-missing branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckAttackRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkAttackRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a"},
						{Name: "b"},
					},
					Attacks: []ast.Attack{
						{Attacker: "a", Target: "b", Line: 1},
					},
				},
			},
		}
		errs := checkAttackRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("MissingTarget", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a"},
					},
					Attacks: []ast.Attack{
						{Attacker: "a", Target: "missing", Line: 2},
					},
				},
			},
		}
		errs := checkAttackRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "target") {
			t.Errorf("expected target error, got %v", errs[0])
		}
	})

	t.Run("MissingAttacker", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "b"},
					},
					Attacks: []ast.Attack{
						{Attacker: "missing", Target: "b", Line: 3},
					},
				},
			},
		}
		errs := checkAttackRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "attacker") {
			t.Errorf("expected attacker error, got %v", errs[0])
		}
	})

	t.Run("BothMissing", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{},
					Attacks: []ast.Attack{
						{Attacker: "missing-a", Target: "missing-b", Line: 4},
					},
				},
			},
		}
		errs := checkAttackRefs(doc)
		if len(errs) != 2 {
			t.Fatalf("expected 2 errors, got %v", errs)
		}
	})
}

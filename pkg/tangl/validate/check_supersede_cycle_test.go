//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckSupersedeCycle — tests checkSupersedeCycle for no-cycle, cycle-detected, and non-pollution of the attack graph
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckSupersedeCycle(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		errs := checkSupersedeCycle(&ast.Document{})
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("NoCycle", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "a", Line: 1}, {Name: "b", Line: 2}},
					Supersedes: []ast.Supersede{
						{Superseder: "a", Superseded: "b"},
					},
				},
			},
		}
		errs := checkSupersedeCycle(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DetectsCycle", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "a", Line: 1}, {Name: "b", Line: 2}},
					Supersedes: []ast.Supersede{
						{Superseder: "a", Superseded: "b"},
						{Superseder: "b", Superseded: "a"},
					},
				},
			},
		}
		errs := checkSupersedeCycle(doc)
		if len(errs) != 1 {
			t.Fatalf("expected exactly 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "supersedes cycle detected") {
			t.Errorf("expected supersedes cycle error, got %v", errs[0])
		}
	})

	t.Run("AttackCycleDoesNotPolluteSupersede", func(t *testing.T) {
		// An Attacks cycle must NOT be reported by the supersede check, and a
		// supersede cycle must NOT be reported by the attack check — the two
		// graphs are independent.
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "a", Line: 1}, {Name: "b", Line: 2}},
					Attacks: []ast.Attack{
						{Attacker: "a", Target: "b"},
						{Attacker: "b", Target: "a"},
					},
				},
			},
		}
		if errs := checkSupersedeCycle(doc); len(errs) != 0 {
			t.Fatalf("attack cycle leaked into supersede check: %v", errs)
		}
		if errs := checkAttackCycle(doc); len(errs) != 1 {
			t.Fatalf("expected attack cycle still detected, got %v", errs)
		}
	})
}

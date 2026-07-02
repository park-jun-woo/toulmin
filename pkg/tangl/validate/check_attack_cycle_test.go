//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckAttackCycle — tests checkAttackCycle for no-cycle and cycle-detected branches across cases
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckAttackCycle(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkAttackCycle(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("NoCycle", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Line: 1},
						{Name: "b", Line: 2},
					},
					Attacks: []ast.Attack{
						{Attacker: "a", Target: "b"},
					},
				},
			},
		}
		errs := checkAttackCycle(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DetectsCycleAcrossMultipleCases", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Line: 1},
						{Name: "b", Line: 2},
					},
					Attacks: []ast.Attack{
						{Attacker: "a", Target: "b"},
						{Attacker: "b", Target: "a"},
					},
				},
				{
					Name: "y",
					Nodes: []ast.Node{
						{Name: "c", Line: 5},
					},
				},
			},
		}
		errs := checkAttackCycle(doc)
		if len(errs) != 1 {
			t.Fatalf("expected exactly 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "cycle detected") {
			t.Errorf("expected cycle detected error, got %v", errs[0])
		}
	})
}

//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseAttacks — tests registerCaseAttacks for empty attacks, ordered attacks, and explicit forces
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRegisterCaseAttacks(t *testing.T) {
	nodes := map[string]nodeInfo{
		"a": {Var: "nA"},
		"b": {Var: "nB"},
		"c": {Var: "nC"},
	}

	t.Run("no attacks emits trailing newline only", func(t *testing.T) {
		var w strings.Builder
		registerCaseAttacks(&w, ast.Case{}, nodes)
		if w.String() != "\n" {
			t.Errorf("got %q, want %q", w.String(), "\n")
		}
	})

	t.Run("multiple attacks emitted in order", func(t *testing.T) {
		var w strings.Builder
		c := ast.Case{
			Attacks: []ast.Attack{
				{Attacker: "a", Target: "b"},
				{Attacker: "b", Target: "c"},
			},
		}
		registerCaseAttacks(&w, c, nodes)
		want := "\tnA.Attacks(nB)\n\tnB.Attacks(nC)\n\n"
		if w.String() != want {
			t.Errorf("got %q, want %q", w.String(), want)
		}
	})

	t.Run("explicit force is passed through, even 1", func(t *testing.T) {
		var w strings.Builder
		half, one := 0.5, 1.0
		c := ast.Case{
			Attacks: []ast.Attack{
				{Attacker: "a", Target: "b", Force: &half},
				{Attacker: "c", Target: "b", Force: &one},
				{Attacker: "b", Target: "c"},
			},
		}
		registerCaseAttacks(&w, c, nodes)
		want := "\tnA.Attacks(nB, 0.5)\n\tnC.Attacks(nB, 1)\n\tnB.Attacks(nC)\n\n"
		if w.String() != want {
			t.Errorf("got %q, want %q", w.String(), want)
		}
	})
}

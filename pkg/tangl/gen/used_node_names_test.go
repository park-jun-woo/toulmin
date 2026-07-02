//ff:func feature=tangl type=codegen control=sequence
//ff:what TestUsedNodeNames — tests usedNodeNames for empty case, attacks-only, and execs-only branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestUsedNodeNames(t *testing.T) {
	t.Run("empty case yields empty map", func(t *testing.T) {
		got := usedNodeNames(ast.Case{})
		if len(got) != 0 {
			t.Errorf("expected empty map, got %v", got)
		}
	})

	t.Run("attacks mark both attacker and target", func(t *testing.T) {
		c := ast.Case{Attacks: []ast.Attack{{Attacker: "a", Target: "b"}}}
		got := usedNodeNames(c)
		if !got["a"] || !got["b"] {
			t.Errorf("expected a and b marked used, got %v", got)
		}
		if len(got) != 2 {
			t.Errorf("expected exactly 2 entries, got %d", len(got))
		}
	})

	t.Run("execs mark node used", func(t *testing.T) {
		c := ast.Case{Execs: []ast.Exec{{Node: "c"}}}
		got := usedNodeNames(c)
		if !got["c"] {
			t.Errorf("expected c marked used, got %v", got)
		}
	})

	t.Run("attacks and execs combined", func(t *testing.T) {
		c := ast.Case{
			Attacks: []ast.Attack{{Attacker: "a", Target: "b"}},
			Execs:   []ast.Exec{{Node: "c"}},
		}
		got := usedNodeNames(c)
		if !got["a"] || !got["b"] || !got["c"] {
			t.Errorf("expected a, b, c all marked used, got %v", got)
		}
	})
}

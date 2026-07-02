//ff:func feature=tangl type=validator control=sequence
//ff:what TestCaseAttackEdges — tests caseAttackEdges for empty-attacks and multiple-attacks branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCaseAttackEdges(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		c := ast.Case{Name: "x"}
		edges := caseAttackEdges(c)
		if len(edges) != 0 {
			t.Fatalf("expected empty edges map, got %v", edges)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		c := ast.Case{
			Name: "x",
			Attacks: []ast.Attack{
				{Attacker: "a", Target: "b"},
				{Attacker: "a", Target: "c"},
				{Attacker: "b", Target: "a"},
			},
		}
		edges := caseAttackEdges(c)
		want := map[string][]string{
			"a": {"b", "c"},
			"b": {"a"},
		}
		if !reflect.DeepEqual(edges, want) {
			t.Fatalf("expected edges=%v, got %v", want, edges)
		}
	})
}

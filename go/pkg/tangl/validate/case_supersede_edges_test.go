//ff:func feature=tangl type=validator control=sequence
//ff:what TestCaseSupersedeEdges — tests caseSupersedeEdges for empty and multi-target superseder mappings
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCaseSupersedeEdges(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		edges := caseSupersedeEdges(ast.Case{})
		if len(edges) != 0 {
			t.Fatalf("expected no edges, got %v", edges)
		}
	})

	t.Run("SupersederToMultipleTargets", func(t *testing.T) {
		c := ast.Case{
			Supersedes: []ast.Supersede{
				{Superseder: "a", Superseded: "b"},
				{Superseder: "a", Superseded: "c"},
			},
		}
		edges := caseSupersedeEdges(c)
		if len(edges["a"]) != 2 || edges["a"][0] != "b" || edges["a"][1] != "c" {
			t.Fatalf("expected a -> [b c], got %v", edges["a"])
		}
	})
}

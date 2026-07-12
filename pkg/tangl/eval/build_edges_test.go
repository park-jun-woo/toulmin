//ff:func feature=tangl type=engine control=sequence
//ff:what TestBuildEdges — tests dependency-edge extraction for rules with and without calls
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildEdges(t *testing.T) {
	rules := []ast.InlineRule{
		{Name: "a", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}},
		{Name: "b", Cond: ast.Call{Fn: ast.Ref{Name: "a"}}},
	}
	edges := buildEdges(rules, map[string]bool{"a": true, "b": true})
	if len(edges) != 2 {
		t.Fatalf("want 2 edge entries, got %d", len(edges))
	}
	if len(edges["a"]) != 0 {
		t.Fatalf("rule a should have no deps, got %v", edges["a"])
	}
	if len(edges["b"]) != 1 || edges["b"][0] != "a" {
		t.Fatalf("rule b should depend on a, got %v", edges["b"])
	}
}

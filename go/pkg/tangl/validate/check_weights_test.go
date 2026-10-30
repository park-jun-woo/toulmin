//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestCheckWeights — valid qualifier/force pass; out-of-range qualifier and force, and a too-fine force, are each reported
package validate

import (
	"math"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckWeights(t *testing.T) {
	ok, bad, fine, nan := 0.5, 1.5, 0.12345, math.NaN()
	doc := &ast.Document{Path: "w.md", Cases: []ast.Case{{
		Name: "c",
		Nodes: []ast.Node{
			{Name: "a", Qualified: &ok},
			{Name: "b"},
			{Name: "q", Qualified: &bad, Line: 3},
		},
		Attacks: []ast.Attack{
			{Attacker: "b", Target: "a"},
			{Attacker: "b", Target: "a", Force: &ok},
			{Attacker: "q", Target: "a", Force: &fine, Line: 5},
			{Attacker: "q", Target: "b", Force: &nan, Line: 6},
		},
	}}}
	errs := checkWeights(doc)
	if len(errs) != 3 {
		t.Fatalf("got %d errors %v, want 3", len(errs), errs)
	}
	for i, want := range []string{`w.md:3: case "c": node "q" qualified 1.5: weight out of range`, `w.md:5: case "c": attack "q" -> "a" force 0.12345: weight has more than 4 decimal places`, `w.md:6: case "c": attack "q" -> "b" force NaN: weight is NaN`} {
		if !strings.HasPrefix(errs[i].Error(), want) {
			t.Errorf("errs[%d] = %q, want prefix %q", i, errs[i], want)
		}
	}
	if errs := checkWeights(&ast.Document{Cases: []ast.Case{{Nodes: []ast.Node{{Qualified: &ok}}, Attacks: []ast.Attack{{Force: &ok}}}}}); len(errs) != 0 {
		t.Errorf("valid weights reported %v", errs)
	}
}

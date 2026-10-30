//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestCheckRuleBindings — rules need exactly one of cond/host; bound nodes may not check or repeat rule properties
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckRuleBindings(t *testing.T) {
	cond := ast.Compare{Field: "x", Op: "is empty"}
	doc := &ast.Document{
		Path: "p.md",
		Rules: []ast.InlineRule{
			{Name: "expr", Line: 1, Cond: cond},
			{Name: "host", Line: 2, Using: &ast.Ref{Name: "f"}},
			{Name: "neither", Line: 3},
			{Name: "both", Line: 4, Cond: cond, Using: &ast.Ref{Name: "f"}},
		},
		Cases: []ast.Case{{Name: "c", Nodes: []ast.Node{
			{Name: "expr", Line: 10},
			{Name: "free", Line: 11, Using: &ast.Ref{Name: "g"}, With: []string{"w"}},
			{Name: "host", Line: 12, Checking: &ast.Ref{Name: "other"}},
			{Name: "expr", Line: 13, With: []string{"w"}},
			{Name: "host", Line: 14, Grounds: []string{"g"}},
			{Name: "expr", Line: 15, Using: &ast.Ref{Name: "h"}},
		}}},
	}
	errs := checkRuleBindings(doc)
	want := []string{
		`p.md:3: rule "neither" must declare exactly one`,
		`p.md:4: rule "both" must declare exactly one`,
		`p.md:12: case "c": checking node "host" conflicts with a rule declaration`,
		`p.md:13: case "c": node "expr" repeats predicate properties`,
		`p.md:14: case "c": node "host" repeats predicate properties`,
		`p.md:15: case "c": node "expr" repeats predicate properties`,
	}
	if len(errs) != len(want) {
		t.Fatalf("errs = %v, want %d", errs, len(want))
	}
	for i, w := range want {
		if !strings.HasPrefix(errs[i].Error(), w) {
			t.Errorf("errs[%d] = %v, want prefix %q", i, errs[i], w)
		}
	}
}

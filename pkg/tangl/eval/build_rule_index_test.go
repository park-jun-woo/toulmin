//ff:func feature=tangl type=engine control=sequence
//ff:what TestBuildRuleIndex — tests the per-document rule lookup tables in document order
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildRuleIndex(t *testing.T) {
	doc := &ast.Document{Rules: []ast.InlineRule{
		{Name: "a", Line: 3},
		{Name: "b", Line: 7},
	}}
	idx := buildRuleIndex(doc)
	if !idx.ruleSet["a"] || !idx.ruleSet["b"] {
		t.Fatalf("ruleSet incomplete: %v", idx.ruleSet)
	}
	if idx.lineOf["a"] != 3 || idx.lineOf["b"] != 7 {
		t.Fatalf("lineOf mismatch: %v", idx.lineOf)
	}
	if idx.byName["b"].Name != "b" {
		t.Fatalf("byName mismatch: %+v", idx.byName)
	}
	if len(idx.names) != 2 || idx.names[0] != "a" || idx.names[1] != "b" {
		t.Fatalf("names should keep document order, got %v", idx.names)
	}
}

//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateCycle — validate catches cycles in attack edges
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateCycle tests that validation catches cycles in attack edges.
func TestValidateCycle(t *testing.T) {
	f := &parser.File{
		Graphs: []parser.GraphDecl{
			{Name: "g", ID: "test:g", Line: 1},
		},
		Bindings: []parser.RuleBinding{
			{Name: "a", Role: "rule", Graph: "g", Func: "fnA", Line: 3},
			{Name: "b", Role: "counter", Graph: "g", Func: "fnB", Line: 4},
			{Name: "c", Role: "except", Graph: "g", Func: "fnC", Line: 5},
		},
		Rules: []parser.InlineRule{
			{Name: "fnA", Expr: parser.Expr{Field: "x", Operator: "is not nil"}, Line: 7},
			{Name: "fnB", Expr: parser.Expr{Field: "y", Operator: "is not nil"}, Line: 8},
			{Name: "fnC", Expr: parser.Expr{Field: "z", Operator: "is not nil"}, Line: 9},
		},
		Attacks: []parser.AttackDecl{
			{Attacker: "a", Target: "b", Graph: "g", Line: 11},
			{Attacker: "b", Target: "c", Graph: "g", Line: 12},
			{Attacker: "c", Target: "a", Graph: "g", Line: 13},
		},
	}

	err := Validate(f)
	if err == nil {
		t.Fatal("expected validation error for cycle")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Errorf("expected cycle error, got: %s", err.Error())
	}
}

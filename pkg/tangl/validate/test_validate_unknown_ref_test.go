//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateUnknownRef — validate catches unknown graph, func, and attack references
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateUnknownRef tests that validation catches unknown references.
func TestValidateUnknownRef(t *testing.T) {
	f := &parser.File{
		Graphs: []parser.GraphDecl{
			{Name: "access", ID: "api:access", Line: 1},
		},
		Bindings: []parser.RuleBinding{
			{Name: "auth", Role: "rule", Graph: "nonexistent", Func: "unknownFunc", Line: 3},
		},
		Attacks: []parser.AttackDecl{
			{Attacker: "ghost", Target: "auth", Graph: "access", Line: 5},
		},
		Evals: []parser.EvalDecl{
			{Name: "result", Graph: "missingGraph", Line: 7},
		},
	}

	err := Validate(f)
	if err == nil {
		t.Fatal("expected validation error for unknown refs")
	}
	msg := err.Error()
	if !strings.Contains(msg, "unknown graph") {
		t.Errorf("expected unknown graph error, got: %s", msg)
	}
	if !strings.Contains(msg, "unknown function") {
		t.Errorf("expected unknown function error, got: %s", msg)
	}
	if !strings.Contains(msg, "unknown binding") {
		t.Errorf("expected unknown binding error, got: %s", msg)
	}
}

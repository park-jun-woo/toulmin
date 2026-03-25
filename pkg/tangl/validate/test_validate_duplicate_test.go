//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateDuplicate — validate catches duplicate names
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateDuplicate tests that validation catches duplicate graph, rule, and binding names.
func TestValidateDuplicate(t *testing.T) {
	f := &parser.File{
		Imports: []parser.Import{
			{Alias: "policy", Package: "github.com/example/pkg", Line: 1},
		},
		Graphs: []parser.GraphDecl{
			{Name: "access", ID: "api:access", Line: 3},
			{Name: "access", ID: "api:access2", Line: 4},
		},
		Rules: []parser.InlineRule{
			{Name: "isAuth", Expr: parser.Expr{Field: "user", Operator: "is not nil"}, Line: 6},
			{Name: "isAuth", Expr: parser.Expr{Field: "ip", Operator: "is not nil"}, Line: 7},
		},
		Bindings: []parser.RuleBinding{
			{Name: "auth", Role: "rule", Graph: "access", Func: "isAuth", Line: 9},
			{Name: "auth", Role: "counter", Graph: "access", Func: "policy.Check", Line: 10},
		},
	}

	err := Validate(f)
	if err == nil {
		t.Fatal("expected validation error for duplicates")
	}
	msg := err.Error()
	if !strings.Contains(msg, "duplicate graph name") {
		t.Errorf("expected duplicate graph error, got: %s", msg)
	}
	if !strings.Contains(msg, "duplicate inline rule name") {
		t.Errorf("expected duplicate rule error, got: %s", msg)
	}
	if !strings.Contains(msg, "duplicate binding name") {
		t.Errorf("expected duplicate binding error, got: %s", msg)
	}
}

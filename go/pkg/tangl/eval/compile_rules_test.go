//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileRules — tests dependency-order rule compilation and error propagation
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCompileRules(t *testing.T) {
	byName := map[string]ast.InlineRule{
		"a": {Name: "a", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}},
		"b": {Name: "b", Cond: ast.Call{Fn: ast.Ref{Name: "a"}}},
	}
	compiled := map[string]evalFunc{}
	dr := docResolver{host: fakeResolver{}, ruleSet: map[string]bool{"a": true, "b": true}, compiled: compiled}
	if err := compileRules([]string{"a", "b"}, byName, dr, compiled); err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(compiled) != 2 {
		t.Fatalf("want 2 compiled rules, got %d", len(compiled))
	}
	if b, _ := compiled["b"](newScope(testCtx(), nil)); !b {
		t.Fatal("rule b should hold via its sibling rule a")
	}
	bad := map[string]ast.InlineRule{"c": {Name: "c", Cond: nil}}
	if err := compileRules([]string{"c"}, bad, dr, map[string]evalFunc{}); err == nil {
		t.Fatal("want compile error for nil condition")
	}
}

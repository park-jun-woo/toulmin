//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompile — tests single-rule compilation into a rule func and the error path
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCompile(t *testing.T) {
	empty := fakeResolver{}

	fn, err := Compile(ast.InlineRule{Name: "r", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}}, empty)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	b, ev := fn(testCtx(), nil)
	if !b || asFact(t, ev).Where != "x" {
		t.Fatalf("compiled rule should hold with an x witness, got %v %+v", b, ev)
	}
	if _, err := Compile(ast.InlineRule{Name: "r", Cond: nil}, empty); err == nil {
		t.Fatal("want compile error for nil condition")
	}
}

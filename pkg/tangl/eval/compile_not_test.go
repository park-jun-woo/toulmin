//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileNot — tests verdict inversion with witness passthrough and an inner error
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCompileNot(t *testing.T) {
	empty := fakeResolver{}

	ef, err := compileNot(ast.Not{Term: ast.Compare{Field: "y", Op: "equals", Value: "2"}}, empty)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	b, ev := ef(newScope(testCtx(), nil))
	if !b || asFact(t, ev).Where != "y" {
		t.Fatalf("not(falsy) should be true with the inner witness, got %v %+v", b, ev)
	}
	ef2, err := compileNot(ast.Not{Term: ast.Compare{Field: "x", Op: "equals", Value: "1"}}, empty)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if b, _ := ef2(newScope(testCtx(), nil)); b {
		t.Fatal("not(truthy) should be false")
	}
	if _, err := compileNot(ast.Not{Term: nil}, empty); err == nil {
		t.Fatal("want compile error for nil inner term")
	}
}

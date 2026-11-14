//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileEither — tests Either first-true compilation and a term compile error
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCompileEither(t *testing.T) {
	empty := fakeResolver{}
	truthy := ast.Compare{Field: "x", Op: "equals", Value: "1"}
	falsy := ast.Compare{Field: "y", Op: "equals", Value: "2"}

	ef, err := compileEither(ast.Either{Terms: []ast.Expr{falsy, truthy}}, empty)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	b, ev := ef(newScope(testCtx(), nil))
	if !b || asFact(t, ev).Where != "x" {
		t.Fatalf("either should stop at the first true term, got %v %+v", b, ev)
	}
	if _, err := compileEither(ast.Either{Terms: []ast.Expr{nil}}, empty); err == nil {
		t.Fatal("want compile error for nil term")
	}
}

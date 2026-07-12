//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileTerms — tests term-list compilation and first-error propagation
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCompileTerms(t *testing.T) {
	empty := fakeResolver{}
	good := ast.Compare{Field: "x", Op: "equals", Value: "1"}

	out, err := compileTerms([]ast.Expr{good, good}, empty)
	if err != nil || len(out) != 2 {
		t.Fatalf("want 2 compiled terms, got %d (%v)", len(out), err)
	}
	if b, _ := out[0](newScope(testCtx(), nil)); !b {
		t.Fatal("compiled term should evaluate true")
	}
	if _, err := compileTerms([]ast.Expr{good, nil}, empty); err == nil {
		t.Fatal("want compile error for nil term")
	}
}

//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileLogic — tests Logic compilation, a term error, and an unknown operator
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCompileLogic(t *testing.T) {
	empty := fakeResolver{}
	truthy := ast.Compare{Field: "x", Op: "equals", Value: "1"}

	ef, err := compileLogic(ast.Logic{Op: "and", Terms: []ast.Expr{truthy, truthy}}, empty)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if b, _ := ef(newScope(testCtx(), nil)); !b {
		t.Fatal("and of two truthy terms should hold")
	}
	if _, err := compileLogic(ast.Logic{Op: "and", Terms: []ast.Expr{nil}}, empty); err == nil {
		t.Fatal("want compile error for nil term")
	}
	if _, err := compileLogic(ast.Logic{Op: "nand", Terms: []ast.Expr{truthy}}, empty); err == nil {
		t.Fatal("want unknown-operator error")
	}
}

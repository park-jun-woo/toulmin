//ff:func feature=tangl type=engine control=sequence
//ff:what runRule — compiles a condition and evaluates it against a ctx
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// runRule compiles cond into a rule func with res and evaluates it against ctx,
// failing the test on a compile error.
func runRule(t *testing.T, cond ast.Expr, res Resolver, ctx toulmin.Context) (bool, any) {
	t.Helper()
	fn, err := Compile(ast.InlineRule{Name: "r", Cond: cond}, res)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return fn(ctx, nil)
}

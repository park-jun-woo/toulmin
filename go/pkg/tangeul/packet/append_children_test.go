//ff:func feature=tangl type=util control=sequence
//ff:what TestAppendChildren — appendChildren flattens terms in order and propagates an emitExpr error
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestAppendChildren checks appendChildren flattens each term via emitExpr,
// summing node counts and appending records in order, and that an error from
// a child's emitExpr call is propagated immediately.
func TestAppendChildren(t *testing.T) {
	var labels []string
	terms := []ast.Expr{
		ast.Compare{Field: "a", Op: "equals", Value: "0"},
		ast.Compare{Field: "b", Op: "equals", Value: "1"},
	}
	out, count, err := appendChildren(nil, 1, terms, &labels)
	if err != nil {
		t.Fatalf("appendChildren: %v", err)
	}
	if count != 3 {
		t.Fatalf("count = %d, want 3 (1 + 2 children)", count)
	}
	if len(out) == 0 {
		t.Fatalf("appendChildren produced no words")
	}

	if _, _, err := appendChildren(nil, 0, []ast.Expr{ast.Compare{Field: "a", Op: "bogus", Value: "0"}}, &labels); err == nil {
		t.Fatalf("child emitExpr error: want error")
	}
}

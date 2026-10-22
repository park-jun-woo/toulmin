//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitExpr — emitExpr handles every ast.Expr kind and rejects unknown Compare/Logic ops
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitExpr exercises every case of emitExpr's type switch — Compare with
// a NUM value, Compare with a Text value, Logic, Not and Either — plus the
// Compare and Logic op-lookup error paths.
func TestEmitExpr(t *testing.T) {
	var labels []string

	if _, _, err := emitExpr(ast.Compare{Field: "n", Op: "equals", Value: "0"}, &labels); err != nil {
		t.Fatalf("Compare NUM: %v", err)
	}
	if _, _, err := emitExpr(ast.Compare{Field: "s", Op: "equals", Value: `"gold"`}, &labels); err != nil {
		t.Fatalf("Compare Text: %v", err)
	}
	if _, _, err := emitExpr(ast.Compare{Field: "n", Op: "bogus", Value: "0"}, &labels); err == nil {
		t.Fatalf("Compare unknown op: want error")
	}

	if _, _, err := emitExpr(ast.Logic{Op: "and", Terms: []ast.Expr{ast.Compare{Field: "n", Op: "equals", Value: "0"}}}, &labels); err != nil {
		t.Fatalf("Logic: %v", err)
	}
	if _, _, err := emitExpr(ast.Logic{Op: "bogus", Terms: []ast.Expr{ast.Compare{Field: "n", Op: "equals", Value: "0"}}}, &labels); err == nil {
		t.Fatalf("Logic unknown op: want error")
	}

	if _, _, err := emitExpr(ast.Not{Term: ast.Compare{Field: "n", Op: "equals", Value: "0"}}, &labels); err != nil {
		t.Fatalf("Not: %v", err)
	}

	if _, _, err := emitExpr(ast.Either{Terms: []ast.Expr{ast.Compare{Field: "n", Op: "equals", Value: "0"}}}, &labels); err != nil {
		t.Fatalf("Either: %v", err)
	}
	many := make([]ast.Expr, exprCountMax+1)
	for i := range many {
		many[i] = ast.Compare{Field: "n", Op: "is empty"}
	}
	if _, _, err := emitExpr(ast.Logic{Op: "and", Terms: many}, &labels); err == nil {
		t.Fatalf("Logic over 13-bit child count: want error")
	}
	if _, _, err := emitExpr(ast.Either{Terms: many}, &labels); err == nil {
		t.Fatalf("Either over 13-bit child count: want error")
	}
	if _, _, err := emitExpr(ast.Exists{Source: ast.Path{Steps: []string{"s"}}, Binder: "b", Cond: ast.Compare{Field: "x", Op: "is empty"}}, &labels); err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if _, _, err := emitExpr(ast.Call{Fn: ast.Ref{Name: "ok"}}, &labels); err != nil {
		t.Fatalf("Call: %v", err)
	}
	if _, _, err := emitExpr(struct{ ast.Compare }{}, &labels); err == nil {
		t.Fatalf("unknown Expr type: want error")
	}
}

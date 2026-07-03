//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitExpr — emitExpr handles every ast.Expr kind and rejects unknown Compare/Logic ops
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
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
}

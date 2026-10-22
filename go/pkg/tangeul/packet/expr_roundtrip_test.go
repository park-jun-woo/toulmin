//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestExprRoundTrip — the EXPR codec is an exact ast.Expr identity over Compare/Logic/Not/Either
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestExprRoundTrip checks the EXPR codec is an exact ast.Expr identity across
// Compare (NUM and Text values, empty-value ops), Logic, Not and Either plus the
// Phase003 nodes — Exists (one-line, nested-logic Cond and nested Exists), the
// operand-carrying CompareX (path lhs, spec lhs, operand rhs, integer rhs and Text
// rhs, and the field-lhs + operand-rhs case), and Call (aliased and local) — the
// node kinds and shapes used by parser/testdata/rules_expr.md's tangl:Rules.
func TestExprRoundTrip(t *testing.T) {
	trees := []ast.Expr{
		ast.Compare{Field: "overdue count", Op: "is at most", Value: "0"},
		ast.Compare{Field: "membership tier", Op: "equals", Value: `"gold"`},
		ast.Compare{Field: "amount", Op: "is empty"},
		ast.Not{Term: ast.Compare{Field: "hold conflict", Op: "equals", Value: `"none"`}},
		ast.Either{Terms: []ast.Expr{
			ast.Compare{Field: "overdue count", Op: "is at most", Value: "0"},
			ast.Compare{Field: "membership tier", Op: "equals", Value: `"gold"`},
		}},
		ast.Logic{Op: "or", Terms: []ast.Expr{
			ast.Either{Terms: []ast.Expr{
				ast.Compare{Field: "overdue count", Op: "is at most", Value: "0"},
				ast.Logic{Op: "and", Terms: []ast.Expr{
					ast.Compare{Field: "membership tier", Op: "equals", Value: `"gold"`},
				}},
			}},
			ast.Compare{Field: "overdue count", Op: "is less than", Value: "2"},
		}},
		// CompareX: path lhs + Text value (trusted author).
		ast.Compare{Op: "is at least", Value: "0.7", Lhs: ast.Path{Steps: []string{"author", "trustScore"}}},
		// CompareX: path lhs + integer value.
		ast.Compare{Op: "equals", Value: "5", Lhs: ast.Path{Steps: []string{"order", "count"}}},
		// CompareX: spec lhs + operand rhs (within budget).
		ast.Compare{Op: "is at most", Lhs: ast.SpecRef{Steps: []string{"discount", "min price"}}, Rhs: ast.Path{Steps: []string{"basePrice"}}},
		// CompareX: field lhs (single-step, Lhs nil) + operand rhs (ipInList inner).
		ast.Compare{Field: "ip", Op: "equals", Rhs: ast.Path{Steps: []string{"clientIp"}}},
		// CompareX: call lhs (value position) + Text value (high nlp).
		ast.Compare{Op: "is at least", Value: "0.7", Lhs: ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}},
		// Call predicate: aliased (spam) and local (meta ok).
		ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}},
		ast.Call{Fn: ast.Ref{Name: "credit ok"}},
		// Exists one-line (hasCoupon).
		ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c", Cond: ast.Compare{Op: "is at most", Lhs: ast.Path{Steps: []string{"c", "minPrice"}}, Rhs: ast.Path{Steps: []string{"basePrice"}}}},
		// Exists with a nested-logic Cond (activePromotion).
		ast.Exists{Source: ast.Path{Steps: []string{"promotions"}}, Binder: "p", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{Op: "equals", Value: "true", Lhs: ast.Path{Steps: []string{"p", "active"}}},
			ast.Compare{Op: "is at least", Value: "10", Lhs: ast.Path{Steps: []string{"p", "discount"}}},
		}}},
		// Nested Exists (exists inside an exists Cond).
		ast.Exists{Source: ast.Path{Steps: []string{"orders"}}, Binder: "o", Cond: ast.Exists{
			Source: ast.Path{Steps: []string{"o", "items"}}, Binder: "i",
			Cond: ast.Compare{Op: "is at least", Value: "1", Lhs: ast.Path{Steps: []string{"i", "qty"}}},
		}},
		// Mixed tree: Logic combining an Exists and a Compare (combo).
		ast.Logic{Op: "or", Terms: []ast.Expr{
			ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c", Cond: ast.Compare{Op: "is at most", Lhs: ast.Path{Steps: []string{"c", "minPrice"}}, Rhs: ast.Path{Steps: []string{"basePrice"}}}},
			ast.Compare{Field: "credit score", Op: "is at least", Value: "700"},
		}},
	}
	for i, in := range trees {
		ws, err := EncodeExpr(ExprPacket{ID: toulmin.TID(i + 1), Cond: in}, testWire)
		if err != nil {
			t.Fatalf("tree %d: EncodeExpr: %v", i, err)
		}
		r, n, err := core.DecodeRecord(ws)
		if err != nil || n != len(ws) || r.Kind != core.KindEdge || r.Codec != 0x61 {
			t.Fatalf("tree %d: not one PREDICATE-DECLARE record: %+v %d %v", i, r, n, err)
		}
		got, err := DecodeExpr(r, testWire)
		if err != nil || got.ID != toulmin.TID(i+1) || !reflect.DeepEqual(got.Cond, in) {
			t.Fatalf("tree %d: DecodeExpr = %#v, %v; want %#v", i, got.Cond, err, in)
		}
	}
}

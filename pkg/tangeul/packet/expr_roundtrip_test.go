//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestExprRoundTrip — the EXPR codec is an exact ast.Expr identity over Compare/Logic/Not/Either
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestExprRoundTrip checks the EXPR codec is an exact ast.Expr identity across
// Compare (NUM and Text values, empty-value ops), Logic, Not and Either — the
// node kinds and shapes used by test/library-renewal.md's tangl:Rules.
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
	}
	for i, in := range trees {
		ws, err := EncodeExpr(ExprPacket{Cond: in})
		if err != nil {
			t.Fatalf("tree %d: EncodeExpr: %v", i, err)
		}
		f, _, _ := word.DecodeFrame(ws)
		if f.Class != word.ClassPayload {
			t.Fatalf("tree %d: frame class = %v, want Payload", i, f.Class)
		}
		got, n, err := DecodeExpr(ws)
		if err != nil || n != len(ws) || !reflect.DeepEqual(got.Cond, in) {
			t.Fatalf("tree %d: DecodeExpr = %#v, %d, %v; want %#v", i, got.Cond, n, err, in)
		}
	}
}

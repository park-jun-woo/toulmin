//ff:func feature=tangl type=util control=sequence
//ff:what TestParseExprSuccess — parseExpr rebuilds a Compare leaf, a Logic node, a Not node and an Either node from hand-built records
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestParseExprSuccess exercises parseExpr's four successful returns —
// Compare, Logic, Not and Either — each rebuilt from hand-built prefix-order
// records.
func TestParseExprSuccess(t *testing.T) {
	labels := []string{"f"}

	pos := 0
	got, count, err := parseExpr(exprCompareRecord(), &pos, labels)
	if err != nil || count != 1 || pos != 6 {
		t.Fatalf("Compare: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	if want := (ast.Compare{Field: "f", Op: "equals", Value: "0"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("Compare = %#v, want %#v", got, want)
	}

	logicRecs := append([]core.Word{exprRecHeader(ExprLogic, 1), core.Word(toulmin.LogicAnd)}, exprCompareRecord()...)
	pos = 0
	got, count, err = parseExpr(logicRecs, &pos, labels)
	if err != nil || count != 2 || pos != len(logicRecs) {
		t.Fatalf("Logic: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	if want := (ast.Logic{Op: "and", Terms: []ast.Expr{ast.Compare{Field: "f", Op: "equals", Value: "0"}}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("Logic = %#v, want %#v", got, want)
	}

	notRecs := append([]core.Word{exprRecHeader(ExprNot, 1)}, exprCompareRecord()...)
	pos = 0
	got, count, err = parseExpr(notRecs, &pos, labels)
	if err != nil || count != 2 || pos != len(notRecs) {
		t.Fatalf("Not: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	if want := (ast.Not{Term: ast.Compare{Field: "f", Op: "equals", Value: "0"}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("Not = %#v, want %#v", got, want)
	}

	eitherRecs := append([]core.Word{exprRecHeader(ExprEither, 1)}, exprCompareRecord()...)
	pos = 0
	got, count, err = parseExpr(eitherRecs, &pos, labels)
	if err != nil || count != 2 || pos != len(eitherRecs) {
		t.Fatalf("Either: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	if want := (ast.Either{Terms: []ast.Expr{ast.Compare{Field: "f", Op: "equals", Value: "0"}}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("Either = %#v, want %#v", got, want)
	}
}

//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestParseExpr — parseExpr dispatches every record kind and rejects an unknown kind
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestParseExpr runs the success and error suites, then round-trips one tree
// of each Phase003 kind (Exists, Call, CompareX) through parseExpr itself so its
// dispatch arms are exercised, and checks an unknown ExprKind is rejected.
func TestParseExpr(t *testing.T) {
	t.Run("success", TestParseExprSuccess)
	t.Run("errors", TestParseExprErrors)
	trees := []ast.Expr{
		ast.Exists{Source: ast.Path{Steps: []string{"s"}}, Binder: "b", Cond: ast.Compare{Field: "x", Op: "is empty"}},
		ast.Call{Fn: ast.Ref{Alias: "a", Name: "ok"}},
		ast.Compare{Op: "equals", Value: "1", Lhs: ast.Path{Steps: []string{"p", "q"}}},
	}
	for i, in := range trees {
		var labels []string
		recs, _, err := emitExpr(in, &labels)
		if err != nil {
			t.Fatalf("%d: emitExpr: %v", i, err)
		}
		pos := 0
		if got, _, err := parseExpr(recs, &pos, labels); err != nil || !reflect.DeepEqual(got, in) {
			t.Fatalf("%d: got %#v %v", i, got, err)
		}
	}
	pos := 0
	if _, _, err := parseExpr([]core.Word{exprRecHeader(ExprKind(7), 0)}, &pos, nil); err == nil {
		t.Fatal("unknown ExprKind: want error")
	}
}

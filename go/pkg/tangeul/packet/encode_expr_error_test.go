//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeExprErrors — EncodeExpr propagates emitExpr errors
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEncodeExprErrors checks that EncodeExpr propagates the error emitExpr
// returns for an unrecognized ast.Compare.Op string.
func TestEncodeExprErrors(t *testing.T) {
	if _, err := EncodeExpr(ExprPacket{Cond: ast.Compare{Field: "f", Op: "bogus", Value: "0"}}, testWire); err == nil {
		t.Fatalf("unknown Compare op: want error")
	}
}

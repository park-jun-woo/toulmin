//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeExprErrors — EncodeExpr propagates emitExpr and ID-width errors
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEncodeExprErrors checks that EncodeExpr propagates the error emitExpr
// returns for an unrecognized ast.Compare.Op string and an Edge Id that does
// not fit the declared width.
func TestEncodeExprErrors(t *testing.T) {
	if _, err := EncodeExpr(ExprPacket{Cond: ast.Compare{Field: "f", Op: "bogus", Value: "0"}}, testWire); err == nil {
		t.Fatalf("unknown Compare op: want error")
	}
	if _, err := EncodeExpr(ExprPacket{ID: 0x10000, Cond: ast.Compare{Field: "f", Op: "is empty"}}, testWire); err == nil {
		t.Fatalf("ID wider than the Edge Id width: want error")
	}
}

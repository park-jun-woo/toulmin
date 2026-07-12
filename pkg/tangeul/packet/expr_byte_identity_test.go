//ff:func feature=tangl type=util control=sequence
//ff:what TestExprByteIdentity — a no-operand Cond's full framed EncodeExpr output is byte-fixed (backward-compat anchor)
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestExprByteIdentity pins the complete framed EncodeExpr output of a Cond that
// uses no Phase003 syntax (a Logic of two legacy kind-0 Compares). The 3-bit
// header widening left every legacy record's bit pattern untouched (bit 13 stays
// zero for kinds 0-3 with small child counts), so this stream is byte-for-byte
// what pre-Phase003 builds emitted — the hard backward-compatibility guarantee.
func TestExprByteIdentity(t *testing.T) {
	cond := ast.Logic{Op: "and", Terms: []ast.Expr{
		ast.Compare{Field: "x", Op: "equals", Value: "1"},
		ast.Compare{Field: "y", Op: "is empty"},
	}}
	want := []word.Word{
		0x1817, 0x0003, 0x0003, 0x0000, 0x0001, 0x7800, 0x0000, 0x0001, 0x7900,
		0x0000, 0x0000, 0x4002, 0x0311, 0x0000, 0x0303, 0x0000, 0x0000, 0x0000,
		0x0001, 0x0000, 0x0301, 0x0001, 0x0001, 0x0002,
	}
	got, err := EncodeExpr(ExprPacket{Cond: cond})
	if err != nil {
		t.Fatalf("EncodeExpr: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("EncodeExpr bytes changed:\n got  %v\n want %v", got, want)
	}
}

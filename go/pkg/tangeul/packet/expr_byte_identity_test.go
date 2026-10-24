//ff:func feature=tangl type=util control=sequence
//ff:what TestExprByteIdentity — a no-operand Cond's full PREDICATE-DECLARE record is byte-fixed
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestExprByteIdentity pins the complete PREDICATE-DECLARE record of a Cond that
// uses no Phase003 syntax (a Logic of two legacy kind-0 Compares) under the GEUL
// t002 layout: TAG WORD + Length extension, LOCAL Edge Id, node count, a label table
// of UTF8-TEXT Values, then the prefix records with the integer value as an
// embedded UNSIGNED-INTEGER Value. Any change to these bytes is a format change.
func TestExprByteIdentity(t *testing.T) {
	cond := ast.Logic{Op: "and", Terms: []ast.Expr{
		ast.Compare{Field: "x", Op: "equals", Value: "1"},
		ast.Compare{Field: "y", Op: "is empty"},
	}}
	want := []core.Word{
		0x1761, 0x0018, // Edge PREDICATE-DECLARE, 1-word Length extension: 24 words
		0x00E0, 0x0001, 0x0003, // LOCAL Edge Id 1, node count 3
		0x0003, 0x08C1, 0x7800, 0x08C1, 0x7900, 0x0840, // labels "x", "y", ""
		0x4002, 0x0311, // Logic and, 2 children
		0x0000, 0x0303, 0x0000, 0x0000, 0x0882, 0x0001, // x equals UNSIGNED-INTEGER 1
		0x0000, 0x0301, 0x0001, 0x0001, 0x0002, // y is empty, Text label 2
	}
	got, err := EncodeExpr(ExprPacket{ID: 1, Cond: cond}, testWire)
	if err != nil {
		t.Fatalf("EncodeExpr: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("EncodeExpr bytes changed:\n got  %v\n want %v", got, want)
	}
}

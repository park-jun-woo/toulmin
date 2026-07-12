//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitCompare — emitCompare emits legacy kind-0 and ExprCompareX golden records and propagates op/operand errors
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitCompare locks emitCompare's two record layouts (§design 3): an
// operand-less Compare keeps the legacy kind-0 [op][field][value] bytes, an
// operand-carrying Compare becomes an ExprCompareX record with an lhs presence
// flag and a ValueOperand-tagged or literal right side. An unknown CompareOp
// and a left or right operand of an unknown type are hard errors.
func TestEmitCompare(t *testing.T) {
	if _, _, err := emitCompare(ast.Compare{Field: "n", Op: "bogus", Value: "7"}, new([]string)); err == nil {
		t.Error("unknown CompareOp: expected error")
	}

	num := word.EncodeNum32(7)
	var labels []string
	recs, n, err := emitCompare(ast.Compare{Field: "n", Op: "equals", Value: "7"}, &labels)
	if err != nil || n != 1 {
		t.Fatalf("legacy: count %d, err %v; want 1, nil", n, err)
	}
	if want := []word.Word{0x0000, 0x0303, 0, 0, num[0], num[1]}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("legacy: recs %v, want %v", recs, want)
	}

	bogus := struct{ ast.Path }{}
	if _, _, err := emitCompare(ast.Compare{Op: "equals", Lhs: bogus}, new([]string)); err == nil {
		t.Error("unknown Lhs operand: expected error")
	}
	if _, _, err := emitCompare(ast.Compare{Op: "equals", Rhs: bogus}, new([]string)); err == nil {
		t.Error("unknown Rhs operand: expected error")
	}

	labels = nil
	x := ast.Compare{Op: "equals", Lhs: ast.Path{Steps: []string{"x"}}, Rhs: ast.SpecRef{Steps: []string{"y"}}}
	recs, n, err = emitCompare(x, &labels)
	if err != nil || n != 1 {
		t.Fatalf("X both operands: count %d, err %v; want 1, nil", n, err)
	}
	if want := []word.Word{0xA000, 0x0303, 0, 1, 0, 1, 1, 2, 1, 1, 2}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("X both operands: recs %v, want %v", recs, want)
	}

	labels = nil
	lhsLit := ast.Compare{Op: "equals", Lhs: ast.Path{Steps: []string{"x"}}, Value: "7"}
	recs, n, err = emitCompare(lhsLit, &labels)
	if err != nil || n != 1 {
		t.Fatalf("X literal right side: count %d, err %v; want 1, nil", n, err)
	}
	if want := []word.Word{0xA000, 0x0303, 0, 1, 0, 1, 1, 0, num[0], num[1]}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("X literal right side: recs %v, want %v", recs, want)
	}
}

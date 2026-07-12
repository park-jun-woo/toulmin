//ff:func feature=tangl type=util control=sequence
//ff:what TestParseCompareX — parseCompareX rebuilds literal and operand-carrying Compares and rejects every malformed record
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseCompareX exercises parseCompareX's successful rebuilds — a literal
// NUM right side without a left operand, and a two-operand Compare via the
// lhs presence flag plus a ValueOperand right side — and every guard: a
// truncated record, an unknown CompareOp, an out-of-range field label, an
// invalid lhs flag, propagated left/right operand errors, a missing value
// kind and an unknown ValueKind.
func TestParseCompareX(t *testing.T) {
	labels := []string{"", "f", "x"}
	eq := word.Word(codebook.CompareEquals)
	zero := word.EncodeNum32(0)

	recs := []word.Word{eq, 1, 0, 0, zero[0], zero[1]}
	pos := 0
	got, count, err := parseCompareX(recs, &pos, labels)
	if err != nil || count != 1 || pos != len(recs) {
		t.Fatalf("literal: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	if want := (ast.Compare{Field: "f", Op: "equals", Value: "0"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("literal = %#v, want %#v", got, want)
	}

	recs = []word.Word{eq, 0, 1, 0, 1, 2, 2, 1, 1, 2}
	pos = 0
	got, count, err = parseCompareX(recs, &pos, labels)
	if err != nil || count != 1 || pos != len(recs) {
		t.Fatalf("operands: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	want := ast.Compare{Op: "equals", Lhs: ast.Path{Steps: []string{"x"}}, Rhs: ast.SpecRef{Steps: []string{"x"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("operands = %#v, want %#v", got, want)
	}

	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 1}, &pos, labels); err == nil {
		t.Error("truncated CompareX record: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{0xFFFF, 1, 0}, &pos, labels); err == nil {
		t.Error("unknown CompareOp: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 9, 0}, &pos, labels); err == nil {
		t.Error("out-of-range field label: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 1, 2}, &pos, labels); err == nil {
		t.Error("invalid lhs presence flag: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 1, 1}, &pos, labels); err == nil {
		t.Error("propagated lhs operand error: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 1, 0}, &pos, labels); err == nil {
		t.Error("missing value kind: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 1, 0, word.Word(ValueOperand)}, &pos, labels); err == nil {
		t.Error("propagated rhs operand error: expected error")
	}
	pos = 0
	if _, _, err := parseCompareX([]word.Word{eq, 1, 0, 9}, &pos, labels); err == nil {
		t.Error("unknown ValueKind: expected error")
	}
}

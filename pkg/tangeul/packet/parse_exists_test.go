//ff:func feature=tangl type=util control=sequence
//ff:what TestParseExists — parseExists rebuilds an ast.Exists from [Quant][source][binder]+Cond and rejects every malformed record
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseExists exercises parseExists's successful rebuild — QuantAny, a
// Path source, a binder label and one Cond child, node count 2 — plus every
// guard: a missing quantifier, an unknown QuantOp, a propagated source
// operand error, a missing binder, an out-of-range binder label, a propagated
// child error and a child count other than one.
func TestParseExists(t *testing.T) {
	labels := []string{"items", "it", "n"}
	quant := word.Word(codebook.QuantAny)
	zero := word.EncodeNum32(0)

	recs := []word.Word{quant, 0, 1, 0, 1, 0x0000, 0x0303, 2, 0, zero[0], zero[1]}
	pos := 0
	got, count, err := parseExists(recs, &pos, labels, 1)
	if err != nil || count != 2 || pos != len(recs) {
		t.Fatalf("Exists: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	want := ast.Exists{Source: ast.Path{Steps: []string{"items"}}, Binder: "it", Cond: ast.Compare{Field: "n", Op: "equals", Value: "0"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Exists = %#v, want %#v", got, want)
	}

	pos = 0
	if _, _, err := parseExists(nil, &pos, labels, 1); err == nil {
		t.Error("missing quantifier: expected error")
	}
	pos = 0
	if _, _, err := parseExists([]word.Word{0xFFFF}, &pos, labels, 1); err == nil {
		t.Error("unknown QuantOp: expected error")
	}
	pos = 0
	if _, _, err := parseExists([]word.Word{quant}, &pos, labels, 1); err == nil {
		t.Error("propagated source operand error: expected error")
	}
	pos = 0
	if _, _, err := parseExists([]word.Word{quant, 0, 0}, &pos, labels, 1); err == nil {
		t.Error("missing binder: expected error")
	}
	pos = 0
	if _, _, err := parseExists([]word.Word{quant, 0, 0, 9}, &pos, labels, 1); err == nil {
		t.Error("out-of-range binder label: expected error")
	}
	pos = 0
	if _, _, err := parseExists([]word.Word{quant, 0, 0, 1}, &pos, labels, 1); err == nil {
		t.Error("propagated child error: expected error")
	}
	pos = 0
	if _, _, err := parseExists([]word.Word{quant, 0, 0, 1}, &pos, labels, 0); err == nil {
		t.Error("zero Cond children: expected error")
	}
}

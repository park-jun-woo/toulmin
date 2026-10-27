//ff:func feature=tangl type=util control=sequence
//ff:what TestParseOperand — parseOperand decodes Path/SpecRef/Call sub-records and rejects missing/truncated/unknown kinds
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestParseOperand exercises every case of parseOperand's kind switch: a Path
// and a SpecRef with their propagated step errors, a Call with its truncation
// and label-range guards, a missing kind word and an unknown OperandKind.
func TestParseOperand(t *testing.T) {
	labels := []string{"a", "nlp", "score"}

	pos := 0
	if _, err := parseOperand(nil, &pos, labels); err == nil {
		t.Error("missing operand kind: expected error")
	}

	pos = 0
	got, err := parseOperand([]core.Word{0, 1, 0}, &pos, labels)
	if err != nil || pos != 3 || !reflect.DeepEqual(got, ast.Path{Steps: []string{"a"}}) {
		t.Fatalf("Path: got %#v, pos %d, err %v", got, pos, err)
	}
	pos = 0
	if _, err := parseOperand([]core.Word{0}, &pos, labels); err == nil {
		t.Error("Path propagated step error: expected error")
	}

	pos = 0
	got, err = parseOperand([]core.Word{1, 1, 0}, &pos, labels)
	if err != nil || pos != 3 || !reflect.DeepEqual(got, ast.SpecRef{Steps: []string{"a"}}) {
		t.Fatalf("SpecRef: got %#v, pos %d, err %v", got, pos, err)
	}
	pos = 0
	if _, err := parseOperand([]core.Word{1, 2, 0}, &pos, labels); err == nil {
		t.Error("SpecRef propagated step error: expected error")
	}

	pos = 0
	got, err = parseOperand([]core.Word{2, 1, 2}, &pos, labels)
	if err != nil || pos != 3 || !reflect.DeepEqual(got, ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}) {
		t.Fatalf("Call: got %#v, pos %d, err %v", got, pos, err)
	}
	pos = 0
	if _, err := parseOperand([]core.Word{2, 1}, &pos, labels); err == nil {
		t.Error("truncated Call operand: expected error")
	}
	pos = 0
	if _, err := parseOperand([]core.Word{2, 9, 9}, &pos, labels); err == nil {
		t.Error("out-of-range Call label: expected error")
	}

	pos = 0
	if _, err := parseOperand([]core.Word{7}, &pos, labels); err == nil {
		t.Error("unknown OperandKind: expected error")
	}
}

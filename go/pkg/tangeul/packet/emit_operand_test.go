//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitOperand — emitOperand encodes Path, SpecRef and Call sub-records and rejects an unknown operand
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitOperand exercises every case of emitOperand's type switch against
// golden words (§design 2): a Path and a SpecRef as stepped sub-records, a
// Call as [OperandCall][alias][name], and a nil operand as the unknown-type
// hard error, sharing one packet-local label table.
func TestEmitOperand(t *testing.T) {
	var labels []string
	recs, err := emitOperand(ast.Path{Steps: []string{"author", "trustScore"}}, &labels)
	if err != nil {
		t.Fatalf("Path: %v", err)
	}
	if want := []core.Word{0, 2, 0, 1}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("Path: recs %v, want %v", recs, want)
	}
	recs, err = emitOperand(ast.SpecRef{Steps: []string{"discount"}}, &labels)
	if err != nil {
		t.Fatalf("SpecRef: %v", err)
	}
	if want := []core.Word{1, 1, 2}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("SpecRef: recs %v, want %v", recs, want)
	}
	recs, err = emitOperand(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, &labels)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if want := []core.Word{2, 3, 4}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("Call: recs %v, want %v", recs, want)
	}
	if _, err := emitOperand(nil, &labels); err == nil {
		t.Error("emitOperand(nil): expected error")
	}
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitExists — emitExists emits the golden ExprExists record [Quant][source][binder]+Cond and propagates operand/child errors
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitExists locks emitExists's ExprExists record layout (§design 3): the
// header word (child count 1), the QuantAny code, the source operand
// sub-record, the interned binder label, then the Cond subtree, with node
// count 1 plus the Cond's. A source operand of an unknown type and a Cond of
// an unknown type are propagated hard errors.
func TestEmitExists(t *testing.T) {
	cond := ast.Compare{Field: "n", Op: "equals", Value: "0"}
	zero := uintVal(0)
	var labels []string
	recs, n, err := emitExists(ast.Exists{Source: ast.Path{Steps: []string{"items"}}, Binder: "it", Cond: cond}, &labels)
	if err != nil || n != 2 {
		t.Fatalf("emitExists: count %d, err %v; want 2, nil", n, err)
	}
	want := []core.Word{0x2001, core.Word(toulmin.QuantAny), 0, 1, 0, 1, 0x0000, 0x0303, 2, 0, zero[0], zero[1]}
	if !reflect.DeepEqual(recs, want) {
		t.Fatalf("emitExists: recs %v, want %v", recs, want)
	}
	if wantLabels := []string{"items", "it", "n"}; !reflect.DeepEqual(labels, wantLabels) {
		t.Fatalf("labels %q, want %q", labels, wantLabels)
	}
	if _, _, err := emitExists(ast.Exists{Binder: "it", Cond: cond}, new([]string)); err == nil {
		t.Error("nil source operand: expected error")
	}
	if _, _, err := emitExists(ast.Exists{Source: ast.Path{Steps: []string{"items"}}, Binder: "it"}, new([]string)); err == nil {
		t.Error("nil Cond child: expected error")
	}
}

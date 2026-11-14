//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestOperandRoundTrip — emitOperand/parseOperand is an exact ast.Operand identity over Path/SpecRef/Call
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestOperandRoundTrip checks the operand codec is an exact ast.Operand identity
// across Path, SpecRef and Call (aliased and local, the local one interning an
// empty alias string), sharing one packet-local label table.
func TestOperandRoundTrip(t *testing.T) {
	operands := []ast.Operand{
		ast.Path{Steps: []string{"author", "trustScore"}},
		ast.Path{Steps: []string{"coupons"}},
		ast.SpecRef{Steps: []string{"discount", "min price"}},
		ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}},
		ast.Call{Fn: ast.Ref{Name: "credit ok"}},
	}
	for i, in := range operands {
		var labels []string
		recs, err := emitOperand(in, &labels)
		if err != nil {
			t.Fatalf("operand %d: emitOperand: %v", i, err)
		}
		pos := 0
		got, err := parseOperand(recs, &pos, labels)
		if err != nil || pos != len(recs) || !reflect.DeepEqual(got, in) {
			t.Fatalf("operand %d: parseOperand = %#v, pos %d, err %v; want %#v", i, got, pos, err, in)
		}
	}
	if _, err := emitOperand(nil, new([]string)); err == nil {
		t.Error("emitOperand(nil): expected error")
	}
}

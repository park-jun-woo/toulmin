//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitCall — emitCall encodes aliased and local predicate calls as golden ExprCall records
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitCall locks emitCall's ExprCall record layout (§design 3): the header
// word (ExprCall, child count 0), the interned alias label index and the
// interned name label index — a local call interning the empty string as its
// alias — with node count 1, sharing one packet-local label table.
func TestEmitCall(t *testing.T) {
	var labels []string
	recs, n, err := emitCall(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, &labels)
	if err != nil || n != 1 {
		t.Fatalf("aliased: count %d, err %v; want 1, nil", n, err)
	}
	if want := []core.Word{0x6000, 0, 1}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("aliased: recs %v, want %v", recs, want)
	}
	recs, n, err = emitCall(ast.Call{Fn: ast.Ref{Name: "credit ok"}}, &labels)
	if err != nil || n != 1 {
		t.Fatalf("local: count %d, err %v; want 1, nil", n, err)
	}
	if want := []core.Word{0x6000, 2, 3}; !reflect.DeepEqual(recs, want) {
		t.Fatalf("local: recs %v, want %v", recs, want)
	}
	if want := []string{"nlp", "score", "", "credit ok"}; !reflect.DeepEqual(labels, want) {
		t.Fatalf("labels %q, want %q", labels, want)
	}
}

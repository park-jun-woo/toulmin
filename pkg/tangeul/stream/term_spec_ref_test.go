//ff:func feature=tangl type=util control=sequence
//ff:what TestTermSpecRef — no sameAs edge yields nil, a resolvable one yields the Ref, an unregistered Func TID errors
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestTermSpecRef covers termSpecRef's three outcomes: no sameAs edge in rel
// (nil, nil), a sameAs edge whose DstTID is a registered Func node (a non-nil
// *ast.Ref built through resolveFuncRef), and a sameAs edge whose DstTID is not
// in d.funcRef, which resolveFuncRef reports as a hard error.
func TestTermSpecRef(t *testing.T) {
	t.Run("no sameAs edge", func(t *testing.T) {
		d := newDec()
		rel := []packet.EdgePacket{{Prop: codebook.PropMeans, DstTID: 5}}
		ref, err := termSpecRef(d, rel)
		if err != nil {
			t.Fatalf("termSpecRef: %v", err)
		}
		if ref != nil {
			t.Errorf("ref = %+v, want nil", ref)
		}
	})

	t.Run("sameAs edge resolves through registered Func node", func(t *testing.T) {
		d := newDec()
		d.funcRef[9] = funcInfo{alias: 2, name: "Temperature"}
		d.nodeName[2] = "sensor"
		rel := []packet.EdgePacket{{Prop: codebook.PropSameAs, DstTID: 9}}
		ref, err := termSpecRef(d, rel)
		if err != nil {
			t.Fatalf("termSpecRef: %v", err)
		}
		if ref == nil {
			t.Fatalf("ref = nil, want a resolved Ref")
		}
		if ref.Alias != "sensor" || ref.Name != "Temperature" {
			t.Errorf("ref = %+v, want {Alias:sensor Name:Temperature}", ref)
		}
	})

	t.Run("sameAs edge to unregistered Func node errors", func(t *testing.T) {
		d := newDec()
		rel := []packet.EdgePacket{{Prop: codebook.PropSameAs, DstTID: 42}}
		ref, err := termSpecRef(d, rel)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if ref != nil {
			t.Errorf("ref = %+v, want nil on error", ref)
		}
	})
}

//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEdgeRoundTrip — Edge packets (FuncRef/TID dst, F_ONCE/F_GATE) round-trip; nil Gate errors
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEdgeRoundTrip exercises Edge packets: a FuncRef dst (do) with F_ONCE, a
// FuncRef dst with F_ONCE|F_GATE carrying a GATE16, a TID dst (attacks), and a
// TID dst (run) — asserting decode(encode(p)) == p and that F_GATE without a Gate
// is rejected.
func TestEdgeRoundTrip(t *testing.T) {
	gate := word.Gate{Op: word.GateAtLeast, Percent: 75}
	packets := []EdgePacket{
		{Prop: codebook.PropDo, Src: 0x0013, DstFunc: &FuncRef{Alias: 0x0002, Verb: 0xC200}, Flags: codebook.EdgeFlagOnce, Span: word.Span{TID: 1, Off: 1, Len: 1}},
		{Prop: codebook.PropDo, Src: 0x0013, DstFunc: &FuncRef{Alias: 0x0002, Verb: 0xC200}, Flags: codebook.EdgeFlagOnce | codebook.EdgeFlagGate, Gate: &gate, Span: word.Span{TID: 1, Off: 2, Len: 2}},
		{Prop: codebook.PropUndo, Src: 0x0013, DstFunc: &FuncRef{Alias: 0x0002, Verb: 0xA000}, Span: word.Span{TID: 1, Off: 3, Len: 3}},
		{Prop: codebook.PropAttacks, Src: 0x0014, DstTID: 0x0013, Span: word.Span{TID: 1, Off: 4, Len: 4}},
		{Prop: codebook.PropRun, Src: 0x0013, DstTID: 0x0020, Span: word.Span{TID: 1, Off: 5, Len: 5}},
	}
	for i, in := range packets {
		ws, err := EncodeEdge(in)
		if err != nil {
			t.Fatalf("packet %d: EncodeEdge: %v", i, err)
		}
		f, _, _ := word.DecodeFrame(ws)
		if f.Class != word.ClassEdge {
			t.Fatalf("packet %d: frame class = %v, want Edge", i, f.Class)
		}
		got, n, err := DecodeEdge(ws)
		if err != nil || n != len(ws) || !reflect.DeepEqual(got, in) {
			t.Fatalf("packet %d: DecodeEdge = %+v, %d, %v; want %+v", i, got, n, err, in)
		}
	}
	if _, err := EncodeEdge(EdgePacket{Prop: codebook.PropDo, DstFunc: &FuncRef{}, Flags: codebook.EdgeFlagGate}); err == nil {
		t.Fatalf("F_GATE without Gate: want error")
	}
}

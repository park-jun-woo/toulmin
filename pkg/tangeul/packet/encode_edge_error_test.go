//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeEdgeErrors — EncodeEdge rejects an unknown PropCode, a nil FuncRef dst, a Block dst and an invalid Gate
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEncodeEdgeErrors exercises EncodeEdge's remaining guards: an unknown
// PropCode (codebook.PropDstKind fails), a FuncRef PropCode with a nil
// DstFunc, a provenance PropCode (Block dst, not encoded in v0.1), and
// F_GATE set with a Gate whose Percent exceeds 100 (word.EncodeGate fails).
func TestEncodeEdgeErrors(t *testing.T) {
	if _, err := EncodeEdge(EdgePacket{Prop: codebook.PropCode(0x0000), Src: 1, Span: word.Span{TID: 1}}); err == nil {
		t.Fatalf("unknown PropCode: want error")
	}
	if _, err := EncodeEdge(EdgePacket{Prop: codebook.PropDo, Src: 1, Span: word.Span{TID: 1}}); err == nil {
		t.Fatalf("nil DstFunc: want error")
	}
	if _, err := EncodeEdge(EdgePacket{Prop: codebook.PropProvenance, Src: 1, Span: word.Span{TID: 1}}); err == nil {
		t.Fatalf("Block dst: want error")
	}
	badGate := word.Gate{Op: word.GateAtLeast, Percent: 101}
	if _, err := EncodeEdge(EdgePacket{Prop: codebook.PropAttacks, Src: 1, DstTID: 2, Flags: codebook.EdgeFlagGate, Gate: &badGate, Span: word.Span{TID: 1}}); err == nil {
		t.Fatalf("invalid Gate percent: want error")
	}
}

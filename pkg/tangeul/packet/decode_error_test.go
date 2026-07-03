//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeErrors — decoders reject wrong class, truncated bodies and slot/length inconsistency
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeErrors checks the self-contained framing guards: a packet decoded as
// the wrong class, a body shorter than its framed length, and a header flag whose
// slot words then do not add up to the framed length are all hard errors rather
// than silent misreads.
func TestDecodeErrors(t *testing.T) {
	nodeWs, _ := EncodeNode(NodePacket{TID: 1, Type: codebook.NodeField, Span: word.Span{TID: 1}})
	edgeWs, _ := EncodeEdge(EdgePacket{Prop: codebook.PropAttacks, Src: 1, DstTID: 2, Span: word.Span{TID: 1}})

	if _, _, err := DecodeEdge(nodeWs); err == nil {
		t.Fatalf("DecodeEdge on a Node packet: want class-mismatch error")
	}
	if _, _, err := DecodeNode(edgeWs); err == nil {
		t.Fatalf("DecodeNode on an Edge packet: want class-mismatch error")
	}
	if _, _, err := DecodeControl(nodeWs); err == nil {
		t.Fatalf("DecodeControl on a Node packet: want class-mismatch error")
	}
	if _, _, err := DecodeExpr(nodeWs); err == nil {
		t.Fatalf("DecodeExpr on a Node packet: want class-mismatch error")
	}
	if _, _, err := DecodeNode(nodeWs[:len(nodeWs)-1]); err == nil {
		t.Fatalf("DecodeNode on a truncated body: want error")
	}

	// Flip F_QUAL into the header without supplying its slot: the slot words no
	// longer add up to the framed length.
	corrupt := append([]word.Word(nil), nodeWs...)
	corrupt[3] |= word.Word(codebook.NodeFlagQual)
	if _, _, err := DecodeNode(corrupt); err == nil {
		t.Fatalf("DecodeNode with an unbacked flag: want length-consistency error")
	}
}

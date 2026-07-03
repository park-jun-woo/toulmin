//ff:func feature=tangl type=util control=sequence
//ff:what TestNodeDescriptiveClaim — the SPAN-object fallback claim carries its kind in the node header
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestNodeDescriptiveClaim checks the type-ladder fallback (§5.2/§5.6): a node
// whose claim could not be structured lands with Subject=Predicate=0, an ObjSpan
// object anchoring the label, and the object kind carried in the header's
// reserved bits so the decoder reads the SPAN object correctly.
func TestNodeDescriptiveClaim(t *testing.T) {
	in := NodePacket{
		TID:   0x0030,
		Type:  codebook.NodeRule,
		Role:  codebook.RoleGeneral,
		Flags: codebook.NodeFlagClaim,
		Claim: &Claim{Kind: ObjSpan, ObjSpan: word.Span{TID: 1, Off: 100, Len: 12}},
		Span:  word.Span{TID: 1, Off: 100, Len: 12},
	}
	ws, err := EncodeNode(in)
	if err != nil {
		t.Fatalf("EncodeNode: %v", err)
	}
	// header word: objKind ObjSpan (=2) sits at bits 8-7 → 0x0100, plus F_CLAIM (b3).
	header := uint16(ws[3])
	if (header>>7)&0x3 != uint16(ObjSpan) {
		t.Fatalf("header objKind bits = %d, want %d (0x%04X)", (header>>7)&0x3, ObjSpan, header)
	}
	got, _, err := DecodeNode(ws)
	if err != nil || !reflect.DeepEqual(got, in) {
		t.Fatalf("DecodeNode = %+v, %v; want %+v", got, err, in)
	}
}

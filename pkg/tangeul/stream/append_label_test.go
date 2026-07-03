//ff:func feature=tangl type=util control=sequence
//ff:what TestAppendLabel — issues a Text payload onto the prefix and appends a label Edge onto the body
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestAppendLabel checks the one reachable path through appendLabel: it always
// builds a well-formed label Edge (PropLabel, DstText kind, no gate), so
// labelEdge/EncodeEdge cannot fail here — the "if err != nil" branch guards a
// codebook.PropDstKind/F_GATE failure that appendLabel's fixed PropLabel +
// zero-Flags construction can never trigger. This test confirms: a Text
// payload for s is appended to *pre carrying its issued TID, and a label Edge
// (src -> that Text TID) is appended to *body with the given span.
func TestAppendLabel(t *testing.T) {
	e := &enc{}
	var pre, body []word.Word
	span := word.Span{TID: 1, Off: 2, Len: 3}

	if err := appendLabel(e, &pre, &body, 7, "hasBalance", span); err != nil {
		t.Fatalf("appendLabel: %v", err)
	}

	wantPre := packet.EncodeText(packet.TextPacket{Bytes: []byte("hasBalance")})
	if len(pre) != len(wantPre) {
		t.Errorf("pre = %v, want %v", pre, wantPre)
	}
	wantBody, err := packet.EncodeEdge(packet.EdgePacket{Prop: codebook.PropLabel, Src: 7, DstTID: 0, Span: span})
	if err != nil {
		t.Fatalf("EncodeEdge reference: %v", err)
	}
	if len(body) != len(wantBody) {
		t.Errorf("body = %v, want %v", body, wantBody)
	}
	if e.textNext != 1 {
		t.Errorf("textNext = %d, want 1", e.textNext)
	}
}

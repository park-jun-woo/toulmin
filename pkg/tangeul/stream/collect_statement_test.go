//ff:func feature=tangl type=util control=sequence
//ff:what TestCollectStatement — prefix/body grouping over empty, prefix-only, same-span, span-break and decode-error inputs
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestCollectStatement exercises collectStatement's two loops: the prefix
// Payload run (zero and nonzero iterations) and the span-bearing body run
// (zero and nonzero iterations, the span-set and span-match branches, the
// differing-span break, and the spanOfRaw decode-error path).
func TestCollectStatement(t *testing.T) {
	// csNodeRaw encodes a minimal valid Node rawPacket at the given TID and span.
	csNodeRaw := func(t *testing.T, tid uint16, span word.Span) rawPacket {
		t.Helper()
		ws, err := packet.EncodeNode(packet.NodePacket{TID: packet.TID(tid), Type: codebook.NodeNamespace, Span: span})
		if err != nil {
			t.Fatalf("EncodeNode: %v", err)
		}
		return rawPacket{Class: word.ClassNode, Words: ws}
	}

	// csEdgeRaw encodes a minimal valid Edge rawPacket at the given span.
	csEdgeRaw := func(t *testing.T, span word.Span) rawPacket {
		t.Helper()
		ws, err := packet.EncodeEdge(packet.EdgePacket{Prop: codebook.PropRequires, Src: 1, DstTID: 2, Span: span})
		if err != nil {
			t.Fatalf("EncodeEdge: %v", err)
		}
		return rawPacket{Class: word.ClassEdge, Words: ws}
	}

	t.Run("empty", func(t *testing.T) {
		prefix, body, next, err := collectStatement(nil, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if prefix != nil || body != nil || next != 0 {
			t.Errorf("got prefix=%v body=%v next=%d, want all zero", prefix, body, next)
		}
	})

	t.Run("prefixOnlyNoBody", func(t *testing.T) {
		packets := []rawPacket{
			{Class: word.ClassPayload, Words: []word.Word{0x2000}},
			{Class: word.ClassPayload, Words: []word.Word{0x2000}},
		}
		prefix, body, next, err := collectStatement(packets, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(prefix) != 2 || body != nil || next != 2 {
			t.Errorf("got prefix=%d body=%v next=%d, want prefix=2 body=nil next=2", len(prefix), body, next)
		}
	})

	t.Run("bodySameSpan", func(t *testing.T) {
		span := word.Span{Off: 0, Len: 4}
		packets := []rawPacket{
			csNodeRaw(t, 1, span),
			csEdgeRaw(t, span),
		}
		prefix, body, next, err := collectStatement(packets, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if prefix != nil || len(body) != 2 || next != 2 {
			t.Errorf("got prefix=%v body=%d next=%d, want prefix=nil body=2 next=2", prefix, len(body), next)
		}
	})

	t.Run("bodyDifferentSpanBreaks", func(t *testing.T) {
		spanA := word.Span{Off: 0, Len: 4}
		spanB := word.Span{Off: 4, Len: 4}
		packets := []rawPacket{
			csNodeRaw(t, 1, spanA),
			csNodeRaw(t, 2, spanB),
		}
		prefix, body, next, err := collectStatement(packets, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if prefix != nil || len(body) != 1 || next != 1 {
			t.Errorf("got prefix=%v body=%d next=%d, want prefix=nil body=1 next=1", prefix, len(body), next)
		}
	})

	t.Run("spanOfRawError", func(t *testing.T) {
		span := word.Span{Off: 0, Len: 4}
		good := csNodeRaw(t, 1, span)
		// Truncate the encoded words so the frame's declared Len runs past the
		// slice, forcing DecodeNode (via spanOfRaw) to error.
		bad := rawPacket{Class: word.ClassNode, Words: good.Words[:len(good.Words)-1]}
		prefix, body, next, err := collectStatement([]rawPacket{bad}, 0)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if prefix != nil || body != nil || next != 0 {
			t.Errorf("got prefix=%v body=%v next=%d on error, want all zero", prefix, body, next)
		}
	})
}

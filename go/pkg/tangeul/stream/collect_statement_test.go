//ff:func feature=tangl type=util control=sequence
//ff:what TestCollectStatement — prefix/body grouping over empty, prefix-only, same-span, span-break and decode-error inputs
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestCollectStatement exercises collectStatement's two loops: the prefix
// Payload run (zero and nonzero iterations) and the span-bearing body run
// (zero and nonzero iterations, the span-set and span-match branches, the
// differing-span break, and the spanOfRaw decode-error path).
func TestCollectStatement(t *testing.T) {
	// csNodeRaw encodes a minimal valid Node rawPacket at the given TID and span.
	csNodeRaw := func(t *testing.T, tid toulmin.TID, span toulmin.Span) rawPacket {
		t.Helper()
		ws, err := toulmin.EncodeNode(toulmin.NodePacket{TID: tid, Type: toulmin.NodeNamespace, Span: span}, testWire)
		if err != nil {
			t.Fatalf("EncodeNode: %v", err)
		}
		return rawOf(t, ws)
	}

	// csEdgeRaw encodes a minimal valid Edge rawPacket at the given span.
	csEdgeRaw := func(t *testing.T, span toulmin.Span) rawPacket {
		t.Helper()
		ws, err := toulmin.EncodeEdge(toulmin.EdgePacket{Prop: toulmin.PropRequires, Src: 1, DstTID: 2, Span: span}, testWire)
		if err != nil {
			t.Fatalf("EncodeEdge: %v", err)
		}
		return rawOf(t, ws)
	}

	t.Run("empty", func(t *testing.T) {
		prefix, body, next, err := collectStatement(newDec(testWire), nil, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if prefix != nil || body != nil || next != 0 {
			t.Errorf("got prefix=%v body=%v next=%d, want all zero", prefix, body, next)
		}
	})

	t.Run("prefixOnlyNoBody", func(t *testing.T) {
		packets := []rawPacket{
			{Class: classPayload},
			{Class: classPayload},
		}
		prefix, body, next, err := collectStatement(newDec(testWire), packets, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(prefix) != 2 || body != nil || next != 2 {
			t.Errorf("got prefix=%d body=%v next=%d, want prefix=2 body=nil next=2", len(prefix), body, next)
		}
	})

	t.Run("bodySameSpan", func(t *testing.T) {
		span := toulmin.Span{Off: 0, Len: 4}
		packets := []rawPacket{
			csNodeRaw(t, 1, span),
			csEdgeRaw(t, span),
		}
		prefix, body, next, err := collectStatement(newDec(testWire), packets, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if prefix != nil || len(body) != 2 || next != 2 {
			t.Errorf("got prefix=%v body=%d next=%d, want prefix=nil body=2 next=2", prefix, len(body), next)
		}
	})

	t.Run("bodyDifferentSpanBreaks", func(t *testing.T) {
		spanA := toulmin.Span{Off: 0, Len: 4}
		spanB := toulmin.Span{Off: 4, Len: 4}
		packets := []rawPacket{
			csNodeRaw(t, 1, spanA),
			csNodeRaw(t, 2, spanB),
		}
		prefix, body, next, err := collectStatement(newDec(testWire), packets, 0)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if prefix != nil || len(body) != 1 || next != 1 {
			t.Errorf("got prefix=%v body=%d next=%d, want prefix=nil body=1 next=1", prefix, len(body), next)
		}
	})

	t.Run("spanOfRawError", func(t *testing.T) {
		span := toulmin.Span{Off: 0, Len: 4}
		good := csNodeRaw(t, 1, span)
		// Drop the last payload word so the SPAN is short, forcing DecodeNode (via
		// spanOfRaw) to error.
		bad := rawPacket{Class: classNode, Rec: good.Rec}
		bad.Rec.Payload = bad.Rec.Payload[:len(bad.Rec.Payload)-1]
		prefix, body, next, err := collectStatement(newDec(testWire), []rawPacket{bad}, 0)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if prefix != nil || body != nil || next != 0 {
			t.Errorf("got prefix=%v body=%v next=%d on error, want all zero", prefix, body, next)
		}
	})
}

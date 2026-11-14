//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeBody — empty body, an appendDecoded error, and a multi-packet success set the shared span once
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestDecodeBody covers decodeBody's loop and its two internal branches: the
// no-span-bearing-packet hard error (zero iterations), an appendDecoded
// failure propagated immediately, and a multi-packet body where the span is
// captured on the first packet and left alone on the rest.
func TestDecodeBody(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		nodes, edges, span, err := decodeBody(newDec(testWire), nil)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if nodes != nil || edges != nil || span != (toulmin.Span{}) {
			t.Errorf("got nodes=%v edges=%v span=%v on error, want zero values", nodes, edges, span)
		}
	})

	t.Run("appendDecodedError", func(t *testing.T) {
		// A payload-class rawPacket hits appendDecoded's default error branch.
		bad := rawPacket{Class: classPayload}
		nodes, edges, span, err := decodeBody(newDec(testWire), []rawPacket{bad})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if nodes != nil || edges != nil || span != (toulmin.Span{}) {
			t.Errorf("got nodes=%v edges=%v span=%v on error, want zero values", nodes, edges, span)
		}
	})

	t.Run("multiPacketKeepsFirstSpan", func(t *testing.T) {
		spanA := toulmin.Span{Off: 0, Len: 4}
		spanB := toulmin.Span{Off: 4, Len: 4}
		nodeWs, err := toulmin.EncodeNode(toulmin.NodePacket{TID: 1, Type: toulmin.NodeNamespace, Span: spanA}, testWire)
		if err != nil {
			t.Fatalf("EncodeNode: %v", err)
		}
		edgeWs, err := toulmin.EncodeEdge(toulmin.EdgePacket{ID: 2, Prop: toulmin.PropRequires, Src: 1, DstTID: 2, Span: spanB}, testWire)
		if err != nil {
			t.Fatalf("EncodeEdge: %v", err)
		}
		body := []rawPacket{
			rawOf(t, nodeWs),
			rawOf(t, edgeWs),
		}
		nodes, edges, span, err := decodeBody(newDec(testWire), body)
		if err != nil {
			t.Fatalf("decodeBody: %v", err)
		}
		if len(nodes) != 1 || len(edges) != 1 {
			t.Errorf("got %d nodes, %d edges, want 1 and 1", len(nodes), len(edges))
		}
		if span != spanA {
			t.Errorf("span = %v, want first packet's span %v", span, spanA)
		}
	})
}

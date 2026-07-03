//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestSpanOfRaw — Node and Edge cases (success and decode error) plus the non-span-bearing default error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestSpanOfRaw drives spanOfRaw directly (white-box) over every switch arm: a
// well-formed Node rawPacket (returns its SPAN), a malformed Node rawPacket
// (packet.DecodeNode's error propagates), the same pair for Edge, and a
// non-span-bearing class (Payload) which falls to the default hard-error arm.
func TestSpanOfRaw(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	ws, err := word.ReadWords(b)
	if err != nil {
		t.Fatalf("ReadWords: %v", err)
	}
	packets, err := readPackets(ws)
	if err != nil {
		t.Fatalf("readPackets: %v", err)
	}
	var goodNode, goodEdge rawPacket
	for _, p := range packets {
		if p.Class == word.ClassNode && goodNode.Words == nil {
			goodNode = p
		}
		if p.Class == word.ClassEdge && goodEdge.Words == nil {
			goodEdge = p
		}
	}
	if goodNode.Words == nil || goodEdge.Words == nil {
		t.Fatalf("fixture stream did not yield both a Node and an Edge packet")
	}

	badNode := rawPacket{Class: word.ClassNode, Words: word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 0})}
	badEdge := rawPacket{Class: word.ClassEdge, Words: word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 0})}
	notSpanBearing := rawPacket{Class: word.ClassPayload, Words: word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 0})}

	tests := []struct {
		name    string
		rp      rawPacket
		wantErr bool
	}{
		{name: "Node success", rp: goodNode, wantErr: false},
		{name: "Node decode error", rp: badNode, wantErr: true},
		{name: "Edge success", rp: goodEdge, wantErr: false},
		{name: "Edge decode error", rp: badEdge, wantErr: true},
		{name: "default: non-span-bearing class", rp: notSpanBearing, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := spanOfRaw(tt.rp)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("spanOfRaw: %v", err)
			}
		})
	}
}

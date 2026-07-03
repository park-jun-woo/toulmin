//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendDecoded — Node/Edge decode success and failure, plus the unexpected-class hard error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestAppendDecoded exercises every case of appendDecoded's Class switch: a
// well-formed Node packet is decoded and appended, a truncated Node packet
// errors before anything is appended, a well-formed Edge packet is decoded and
// appended, a truncated Edge packet errors, and a Payload-class rawPacket (not
// span-bearing in the body) hits the default hard error.
func TestAppendDecoded(t *testing.T) {
	nodeWords, err := packet.EncodeNode(packet.NodePacket{Role: 0, Span: word.Span{TID: 1, Off: 2, Len: 3}})
	if err != nil {
		t.Fatalf("EncodeNode fixture: %v", err)
	}
	edgeWords, err := packet.EncodeEdge(packet.EdgePacket{Prop: 0x0101, Span: word.Span{TID: 4, Off: 5, Len: 6}})
	if err != nil {
		t.Fatalf("EncodeEdge fixture: %v", err)
	}

	cases := []struct {
		name    string
		rp      rawPacket
		wantErr bool
		wantN   int
		wantE   int
	}{
		{"node ok", rawPacket{Class: word.ClassNode, Words: nodeWords}, false, 1, 0},
		{"node truncated", rawPacket{Class: word.ClassNode, Words: nil}, true, 0, 0},
		{"edge ok", rawPacket{Class: word.ClassEdge, Words: edgeWords}, false, 0, 1},
		{"edge truncated", rawPacket{Class: word.ClassEdge, Words: nil}, true, 0, 0},
		{"unexpected class", rawPacket{Class: word.ClassPayload, Words: []word.Word{0}}, true, 0, 0},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		var nodes []packet.NodePacket
		var edges []packet.EdgePacket
		span, err := appendDecoded(&nodes, &edges, c.rp)
		if c.wantErr {
			if err == nil {
				t.Fatalf("want error, got nil (span %+v)", span)
			}
			return
		}
		if err != nil {
			t.Fatalf("appendDecoded: %v", err)
		}
		if len(nodes) != c.wantN || len(edges) != c.wantE {
			t.Errorf("nodes=%d edges=%d, want nodes=%d edges=%d", len(nodes), len(edges), c.wantN, c.wantE)
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
}

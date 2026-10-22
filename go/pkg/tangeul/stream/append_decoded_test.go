//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendDecoded — NODE and LINK records decode into their lists; bad, repeated-id and span-less records fail
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestAppendDecoded decodes a NODE and a LINK, then checks malformed records, an
// Edge Id seen twice, and a payload-class record.
func TestAppendDecoded(t *testing.T) {
	nw, _ := toulmin.EncodeNode(toulmin.NodePacket{TID: 1, Type: toulmin.NodeCase, Span: toulmin.Span{Off: 2, Len: 3}}, testWire)
	ew, _ := toulmin.EncodeEdge(toulmin.EdgePacket{ID: 2, Prop: toulmin.PropAttacks, Span: toulmin.Span{Off: 5, Len: 6}}, testWire)
	d := newDec(testWire)
	var nodes []toulmin.NodePacket
	var edges []toulmin.EdgePacket
	if s, err := appendDecoded(d, &nodes, &edges, rawOf(t, nw)); err != nil || s.Off != 2 || len(nodes) != 1 {
		t.Fatalf("node: %+v %v", s, err)
	}
	if s, err := appendDecoded(d, &nodes, &edges, rawOf(t, ew)); err != nil || s.Off != 5 || len(edges) != 1 {
		t.Fatalf("edge: %+v %v", s, err)
	}
	bad := []rawPacket{
		rawOf(t, nw),
		rawOf(t, ew),
		{Class: classNode, Rec: core.Record{Kind: core.KindEdge, Codec: toulmin.CodecCaseDeclare}},
		{Class: classLink, Rec: core.Record{Kind: core.KindEdge, Codec: toulmin.CodecAttacks}},
		{Class: classPayload},
	}
	for i, rp := range bad {
		if _, err := appendDecoded(d, &nodes, &edges, rp); err == nil {
			t.Fatalf("bad %d: want error", i)
		}
	}
}

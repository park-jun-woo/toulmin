//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestSpanOfRaw — spanOfRaw returns NODE/LINK spans and rejects malformed and span-less records
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestSpanOfRaw drives every arm of spanOfRaw: a real NODE and LINK from an
// encoded document (success), malformed ones (decode errors propagate), and a
// span-less class (default error).
func TestSpanOfRaw(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	d, packets := docPackets(t, b)
	var node, link *rawPacket
	for i := range packets {
		if packets[i].Class == classNode && node == nil {
			node = &packets[i]
		}
		if packets[i].Class == classLink && link == nil {
			link = &packets[i]
		}
	}
	if node == nil || link == nil {
		t.Fatal("fixture has no NODE or LINK record")
	}
	good := []rawPacket{*node, *link}
	for _, rp := range good {
		if _, err := spanOfRaw(d, rp); err != nil {
			t.Fatalf("spanOfRaw: %v", err)
		}
	}
	bad := []rawPacket{
		{Class: classNode, Rec: core.Record{Kind: core.KindEdge, Codec: toulmin.CodecCaseDeclare}},
		{Class: classLink, Rec: core.Record{Kind: core.KindEdge, Codec: toulmin.CodecAttacks}},
		{Class: classPayload},
	}
	for i, rp := range bad {
		if _, err := spanOfRaw(d, rp); err == nil {
			t.Fatalf("bad %d: want error", i)
		}
	}
}

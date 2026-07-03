//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestSplitEdges — label edges and relation edges partition correctly, including the empty and single-kind cases
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestSplitEdges drives splitEdges over an empty slice (the loop runs zero
// times), an all-label slice, an all-relation slice and a mixed slice, checking
// that every Prop==codebook.PropLabel edge lands in labelEdges (in order) and
// every other edge lands in rel (in order).
func TestSplitEdges(t *testing.T) {
	label1 := packet.EdgePacket{Prop: codebook.PropLabel, DstTID: 1}
	label2 := packet.EdgePacket{Prop: codebook.PropLabel, DstTID: 2}
	req := packet.EdgePacket{Prop: codebook.PropRequires, DstTID: 3}
	means := packet.EdgePacket{Prop: codebook.PropMeans, DstTID: 4}

	cases := []struct {
		name      string
		edges     []packet.EdgePacket
		wantLabel []packet.EdgePacket
		wantRel   []packet.EdgePacket
	}{
		{name: "empty", edges: nil, wantLabel: nil, wantRel: nil},
		{name: "all label", edges: []packet.EdgePacket{label1, label2}, wantLabel: []packet.EdgePacket{label1, label2}, wantRel: nil},
		{name: "all relation", edges: []packet.EdgePacket{req, means}, wantLabel: nil, wantRel: []packet.EdgePacket{req, means}},
		{name: "mixed", edges: []packet.EdgePacket{req, label1, means, label2}, wantLabel: []packet.EdgePacket{label1, label2}, wantRel: []packet.EdgePacket{req, means}},
	}
	for _, c := range cases {
		gotLabel, gotRel := splitEdges(c.edges)
		if !reflect.DeepEqual(gotLabel, c.wantLabel) {
			t.Errorf("%s: labelEdges = %+v, want %+v", c.name, gotLabel, c.wantLabel)
		}
		if !reflect.DeepEqual(gotRel, c.wantRel) {
			t.Errorf("%s: rel = %+v, want %+v", c.name, gotRel, c.wantRel)
		}
	}
}

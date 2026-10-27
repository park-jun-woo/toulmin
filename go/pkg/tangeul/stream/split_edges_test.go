//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestSplitEdges — label edges and relation edges partition correctly, including the empty and single-kind cases
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestSplitEdges drives splitEdges over an empty slice (the loop runs zero
// times), an all-label slice, an all-relation slice and a mixed slice, checking
// that every Prop==toulmin.PropLabel edge lands in labelEdges (in order) and
// every other edge lands in rel (in order).
func TestSplitEdges(t *testing.T) {
	label1 := toulmin.EdgePacket{Prop: toulmin.PropLabel, DstTID: 1}
	label2 := toulmin.EdgePacket{Prop: toulmin.PropLabel, DstTID: 2}
	req := toulmin.EdgePacket{Prop: toulmin.PropRequires, DstTID: 3}
	means := toulmin.EdgePacket{Prop: toulmin.PropMeans, DstTID: 4}

	cases := []struct {
		name      string
		edges     []toulmin.EdgePacket
		wantLabel []toulmin.EdgePacket
		wantRel   []toulmin.EdgePacket
	}{
		{name: "empty", edges: nil, wantLabel: nil, wantRel: nil},
		{name: "all label", edges: []toulmin.EdgePacket{label1, label2}, wantLabel: []toulmin.EdgePacket{label1, label2}, wantRel: nil},
		{name: "all relation", edges: []toulmin.EdgePacket{req, means}, wantLabel: nil, wantRel: []toulmin.EdgePacket{req, means}},
		{name: "mixed", edges: []toulmin.EdgePacket{req, label1, means, label2}, wantLabel: []toulmin.EdgePacket{label1, label2}, wantRel: []toulmin.EdgePacket{req, means}},
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

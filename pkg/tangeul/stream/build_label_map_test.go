//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestBuildLabelMap — buildLabelMap groups label edges by src TID in order, erroring on an unknown Text
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestBuildLabelMap drives buildLabelMap's loop through zero iterations (an
// empty edge list), several iterations that group multiple src TIDs and
// preserve per-src emission order, and an edge whose DstTID has no registered
// Text (a hard error).
func TestBuildLabelMap(t *testing.T) {
	cases := []struct {
		name    string
		texts   map[uint16]string
		edges   []packet.EdgePacket
		wantErr bool
		want    map[uint16][]string
	}{
		{
			name:  "zero edges",
			texts: map[uint16]string{},
			edges: nil,
			want:  map[uint16][]string{},
		},
		{
			name:  "grouped, order preserved",
			texts: map[uint16]string{10: "a", 11: "b", 12: "c"},
			edges: []packet.EdgePacket{
				{Src: 1, DstTID: 10},
				{Src: 2, DstTID: 11},
				{Src: 1, DstTID: 12},
			},
			want: map[uint16][]string{1: {"a", "c"}, 2: {"b"}},
		},
		{
			name:    "unknown text",
			texts:   map[uint16]string{},
			edges:   []packet.EdgePacket{{Src: 1, DstTID: 99}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		d := newDec()
		d.texts = c.texts
		got, err := buildLabelMap(d, c.edges)
		if c.wantErr && err == nil {
			t.Errorf("%s: want error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Errorf("%s: buildLabelMap: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}

//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFindRel — an empty slice, a matching predicate, and no match after a full scan
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestFindRel covers findRel's loop with zero iterations, a match found
// partway through, and a full scan that finds nothing.
func TestFindRel(t *testing.T) {
	cases := []struct {
		name   string
		rel    []packet.EdgePacket
		prop   codebook.PropCode
		want   packet.EdgePacket
		wantOK bool
	}{
		{name: "empty", rel: nil, prop: codebook.PropDo, wantOK: false},
		{
			name: "matchSecond",
			rel: []packet.EdgePacket{
				{Prop: codebook.PropRequires},
				{Prop: codebook.PropDo, Src: 3},
			},
			prop:   codebook.PropDo,
			want:   packet.EdgePacket{Prop: codebook.PropDo, Src: 3},
			wantOK: true,
		},
		{
			name: "noMatch",
			rel: []packet.EdgePacket{
				{Prop: codebook.PropRequires},
				{Prop: codebook.PropUndo},
			},
			prop:   codebook.PropDo,
			wantOK: false,
		},
	}
	for _, c := range cases {
		got, ok := findRel(c.rel, c.prop)
		if ok != c.wantOK {
			t.Errorf("%s: ok = %v, want %v", c.name, ok, c.wantOK)
		}
		if ok && got != c.want {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}

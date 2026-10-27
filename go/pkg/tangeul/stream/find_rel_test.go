//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestFindRel — an empty slice, a matching predicate, and no match after a full scan
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestFindRel covers findRel's loop with zero iterations, a match found
// partway through, and a full scan that finds nothing.
func TestFindRel(t *testing.T) {
	cases := []struct {
		name   string
		rel    []toulmin.EdgePacket
		prop   toulmin.PropCode
		want   toulmin.EdgePacket
		wantOK bool
	}{
		{name: "empty", rel: nil, prop: toulmin.PropDo, wantOK: false},
		{
			name: "matchSecond",
			rel: []toulmin.EdgePacket{
				{Prop: toulmin.PropRequires},
				{Prop: toulmin.PropDo, Src: 3},
			},
			prop:   toulmin.PropDo,
			want:   toulmin.EdgePacket{Prop: toulmin.PropDo, Src: 3},
			wantOK: true,
		},
		{
			name: "noMatch",
			rel: []toulmin.EdgePacket{
				{Prop: toulmin.PropRequires},
				{Prop: toulmin.PropUndo},
			},
			prop:   toulmin.PropDo,
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

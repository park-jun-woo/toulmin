//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestBuildTrigger — buildTrigger sets Event or Interval by kind and resolves Until from its edge
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestBuildTrigger covers both InternalKind arms (OnEvent sets Event, anything
// else sets Interval) crossed with both arms of the until-edge lookup (present
// vs absent), so each independent branch is hit at least once.
func TestBuildTrigger(t *testing.T) {
	cases := []struct {
		name string
		role codebook.RoleCode
		rel  []packet.EdgePacket
		want ast.Internal
	}{
		{
			name: "on-event with until",
			role: codebook.RoleCode(ast.OnEvent),
			rel:  []packet.EdgePacket{{Prop: codebook.PropUntil, DstTID: 2}},
			want: ast.Internal{Kind: ast.OnEvent, Event: "sensor.buttonPress", Until: "stop case"},
		},
		{
			name: "every-tick without until",
			role: codebook.RoleCode(ast.EveryTick),
			rel:  nil,
			want: ast.Internal{Kind: ast.EveryTick, Interval: "30s"},
		},
	}
	for _, c := range cases {
		d := newDec()
		d.nodeName[2] = "stop case"
		name := c.want.Event
		if name == "" {
			name = c.want.Interval
		}
		ref, err := buildTrigger(d, 1, name, c.role, c.rel)
		if err != nil {
			t.Fatalf("%s: buildTrigger: %v", c.name, err)
		}
		if len(d.ast.Internals) != 1 || !reflect.DeepEqual(d.ast.Internals[0], c.want) {
			t.Errorf("%s: Internals = %+v, want [%+v]", c.name, d.ast.Internals, c.want)
		}
		if d.trByTID[1] != 0 {
			t.Errorf("%s: trByTID[1] = %d, want 0", c.name, d.trByTID[1])
		}
		wantRef := doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0}}
		if !reflect.DeepEqual(ref, wantRef) {
			t.Errorf("%s: StmtRef = %+v, want %+v", c.name, ref, wantRef)
		}
	}
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestBuildAttack — an ATTACKS LINK rebuilds attacker/target names, F_FORCE becomes an explicit Force, no F_FORCE stays nil
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestBuildAttack(t *testing.T) {
	d := &dec{nodeName: map[toulmin.TID]string{7: "block ip", 8: "authenticate"}}
	a := buildAttack(d, toulmin.EdgePacket{Prop: toulmin.PropAttacks, Src: 7, DstTID: 8})
	if a.Attacker != "block ip" || a.Target != "authenticate" || a.Force != nil {
		t.Fatalf("omitted: got %+v", a)
	}
	one := toulmin.VFix(10000)
	a = buildAttack(d, toulmin.EdgePacket{Prop: toulmin.PropAttacks, Src: 7, DstTID: 8, Flags: toulmin.EdgeFlagForce, Force: &one})
	if a.Force == nil || *a.Force != 1 {
		t.Fatalf("explicit 1 must stay explicit: got %+v", a)
	}
}

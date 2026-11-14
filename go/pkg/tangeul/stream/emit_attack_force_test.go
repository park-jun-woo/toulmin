//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitAttackForce — explicit force sets F_FORCE + value on the wire, omitted force leaves it off, out-of-range force errors
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitAttackForce encodes an explicit-force attack and decodes the LINK
// back: F_FORCE must be set and carry 0.5. The omitted form must decode with
// no F_FORCE and no Force. A force VFix cannot represent (1.5) is a hard error.
func TestEmitAttackForce(t *testing.T) {
	e := testEnc(&tidTable{nodeByCase: map[string]toulmin.TID{nodeKey(0, "a"): 7, nodeKey(0, "b"): 8}})
	ref := doc.StmtRef{Sec: 5, Path: []int{0, 2, 0}}
	decode := func(words []core.Word) toulmin.EdgePacket {
		t.Helper()
		rec, _, err := core.DecodeRecord(words)
		if err != nil {
			t.Fatalf("DecodeRecord: %v", err)
		}
		ep, err := toulmin.DecodeEdge(rec, testWire)
		if err != nil {
			t.Fatalf("DecodeEdge: %v", err)
		}
		return ep
	}
	half := 0.5
	words, err := emitAttack(e, ref, ast.Attack{Attacker: "a", Target: "b", Force: &half}, toulmin.Span{})
	if err != nil {
		t.Fatalf("emitAttack(0.5): %v", err)
	}
	ep := decode(words)
	if ep.Flags&toulmin.EdgeFlagForce == 0 || ep.Force == nil || ep.Force.Float() != 0.5 || ep.Src != 7 || ep.DstTID != 8 {
		t.Fatalf("explicit force decoded as %+v", ep)
	}
	words, err = emitAttack(e, ref, ast.Attack{Attacker: "a", Target: "b"}, toulmin.Span{})
	if err != nil {
		t.Fatalf("emitAttack(omitted): %v", err)
	}
	if ep := decode(words); ep.Flags&toulmin.EdgeFlagForce != 0 || ep.Force != nil {
		t.Fatalf("omitted force decoded as %+v", ep)
	}
	big := 1.5
	if _, err := emitAttack(e, ref, ast.Attack{Attacker: "a", Target: "b", Force: &big}, toulmin.Span{}); err == nil {
		t.Fatal("force 1.5 must be a hard error")
	}
}

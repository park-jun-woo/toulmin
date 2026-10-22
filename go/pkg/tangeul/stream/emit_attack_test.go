//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitAttack — undeclared attacker, undeclared target, and a successful defeat edge
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitAttack drives emitAttack's two undeclared-node hard errors (attacker,
// then target) and the success path where both resolve through nodeByCase.
func TestEmitAttack(t *testing.T) {
	e := testEnc(&tidTable{
		nodeByCase: map[string]toulmin.TID{
			nodeKey(0, "a"): 7,
			nodeKey(0, "b"): 8,
		},
	})
	ref := doc.StmtRef{Sec: 5, Path: []int{0, 2, 0}}
	span := toulmin.Span{Off: 1, Len: 2}

	cases := []struct {
		name    string
		a       ast.Attack
		wantErr bool
	}{
		{"undeclared attacker", ast.Attack{Attacker: "missing", Target: "b"}, true},
		{"undeclared target", ast.Attack{Attacker: "a", Target: "missing"}, true},
		{"success", ast.Attack{Attacker: "a", Target: "b"}, false},
	}
	for _, c := range cases {
		got, err := emitAttack(e, ref, c.a, span)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: emitAttack: %v", c.name, err)
		}
		if len(got) == 0 {
			t.Errorf("%s: expected non-empty words", c.name)
		}
	}

	atkTT := nodeTT(doc.StmtRef{}, 0)
	atkTT.nodeByCase[nodeKey(0, "a")], atkTT.nodeByCase[nodeKey(0, "b")] = 1, 2
	failsAtEachID(t, 1, atkTT, func(e *enc) error {
		_, err := emitAttack(e, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 2, 0}}, ast.Attack{Attacker: "a", Target: "b"}, toulmin.Span{})
		return err
	})
}

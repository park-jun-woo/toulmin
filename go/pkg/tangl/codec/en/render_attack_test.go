//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderAttack — renderAttack emits the omitted-force `don't` form and the explicit `attacks ... with force` form
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderAttack checks both defeat edge forms: a nil Force renders the
// omitted form, and a set Force — 0.5, and 1 which stays explicit — renders
// "`attacker` attacks `target` with force F" with the shortest decimal.
func TestRenderAttack(t *testing.T) {
	got := renderAttack(ast.Attack{Target: "authenticate", Attacker: "block ip"})
	if got != "don't `authenticate` when `block ip`" {
		t.Errorf("renderAttack = %q", got)
	}
	half, one := 0.5, 1.0
	got = renderAttack(ast.Attack{Target: "authenticate", Attacker: "block ip", Force: &half})
	if got != "`block ip` attacks `authenticate` with force 0.5" {
		t.Errorf("renderAttack(0.5) = %q", got)
	}
	got = renderAttack(ast.Attack{Target: "a", Attacker: "b", Force: &one})
	if got != "`b` attacks `a` with force 1" {
		t.Errorf("renderAttack(1) = %q", got)
	}
}

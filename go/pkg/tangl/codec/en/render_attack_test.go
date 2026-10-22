//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderAttack — renderAttack emits `don't `target` when `attacker``
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderAttack checks the defeat edge form.
func TestRenderAttack(t *testing.T) {
	got := renderAttack(ast.Attack{Target: "authenticate", Attacker: "block ip"})
	if got != "don't `authenticate` when `block ip`" {
		t.Errorf("renderAttack = %q", got)
	}
}

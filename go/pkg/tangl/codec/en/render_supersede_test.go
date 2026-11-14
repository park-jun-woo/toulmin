//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSupersede — renderSupersede emits `` `superseder` supersedes `superseded` ``
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderSupersede checks the crisp-exclusion edge form.
func TestRenderSupersede(t *testing.T) {
	got := renderSupersede(ast.Supersede{Superseder: "overdrawn", Superseded: "pending"})
	if got != "`overdrawn` supersedes `pending`" {
		t.Errorf("renderSupersede = %q", got)
	}
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRef — renderRef emits qualified and bare references
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderRef checks both the qualified `alias`.`name` and bare `name` forms.
func TestRenderRef(t *testing.T) {
	if got := renderRef(ast.Ref{Alias: "arm", Name: "placeCup"}); got != "`arm`.`placeCup`" {
		t.Errorf("qualified = %q", got)
	}
	if got := renderRef(ast.Ref{Name: "fn"}); got != "`fn`" {
		t.Errorf("bare = %q", got)
	}
}

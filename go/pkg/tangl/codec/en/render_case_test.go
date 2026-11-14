//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderCase — renderCase emits `in case of `X``
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCase checks the case declaration head form.
func TestRenderCase(t *testing.T) {
	if got := renderCase(ast.Case{Name: "can access"}); got != "in case of `can access`" {
		t.Errorf("renderCase = %q", got)
	}
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSee — renderSee emits `see `alias` from `package``
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderSee checks the see declaration form.
func TestRenderSee(t *testing.T) {
	got := renderSee(ast.See{Alias: "arm", Package: "robot/arm"})
	if got != "see `arm` from `robot/arm`" {
		t.Errorf("renderSee = %q", got)
	}
}

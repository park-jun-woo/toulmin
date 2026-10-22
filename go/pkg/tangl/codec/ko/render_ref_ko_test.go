//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRefKo — renderRefKo emits `alias`.`name` when qualified and bare `name` otherwise
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderRefKo checks both the qualified and bare-local reference forms.
func TestRenderRefKo(t *testing.T) {
	if got := renderRefKo(ast.Ref{Alias: "센서", Name: "isOrdered"}); got != "`센서`.`isOrdered`" {
		t.Errorf("qualified ref = %q", got)
	}
	if got := renderRefKo(ast.Ref{Name: "authenticate"}); got != "`authenticate`" {
		t.Errorf("bare ref = %q", got)
	}
}

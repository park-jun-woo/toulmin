//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderField — renderField emits `has `f` as <Type>`
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderField checks the struct field form with a built-in type.
func TestRenderField(t *testing.T) {
	got := renderField(ast.Field{Name: "annual income", Type: "Currency"})
	if got != "has `annual income` as Currency" {
		t.Errorf("renderField = %q", got)
	}
}

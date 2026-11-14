//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRequire — renderRequire covers the bare and `as <Type>` forms
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderRequire checks the required-field form with and without a type.
func TestRenderRequire(t *testing.T) {
	if got := renderRequire(ast.Require{Field: "user"}); got != "`user` is required" {
		t.Errorf("bare = %q", got)
	}
	if got := renderRequire(ast.Require{Field: "age", Type: "Integer"}); got != "`age` is required as Integer" {
		t.Errorf("typed = %q", got)
	}
}

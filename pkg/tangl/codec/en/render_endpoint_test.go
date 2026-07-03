//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderEndpoint — renderEndpoint emits `provides `X``
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderEndpoint checks the endpoint declaration head form.
func TestRenderEndpoint(t *testing.T) {
	if got := renderEndpoint(ast.Endpoint{Name: "transfer"}); got != "provides `transfer`" {
		t.Errorf("renderEndpoint = %q", got)
	}
}

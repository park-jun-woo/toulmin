//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderEndpointKo — renderEndpointKo emits `` `X` 를 한다면 `` with 받침-selected 을/를
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderEndpointKo checks the endpoint head on an open name (이체→를) and a 받침
// name (접근→을).
func TestRenderEndpointKo(t *testing.T) {
	cases := []struct {
		e    ast.Endpoint
		want string
	}{
		{ast.Endpoint{Name: "이체"}, "`이체` 를 한다면"},
		{ast.Endpoint{Name: "접근"}, "`접근` 을 한다면"},
	}
	for _, c := range cases {
		if got := renderEndpointKo(c.e); got != c.want {
			t.Errorf("renderEndpointKo(%+v) = %q, want %q", c.e, got, c.want)
		}
	}
}

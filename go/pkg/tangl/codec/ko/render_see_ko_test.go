//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderSeeKo — renderSeeKo emits `` `p` 는 `path` 에서 참조한다 `` with 받침-selected 은/는
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderSeeKo checks the 참조 form and its 은/는 allomorph on a 받침 alias (팔→은)
// and an open alias (센서→는).
func TestRenderSeeKo(t *testing.T) {
	cases := []struct {
		s    ast.See
		want string
	}{
		{ast.See{Alias: "팔", Package: "robot/arm"}, "`팔` 은 `robot/arm` 에서 참조한다"},
		{ast.See{Alias: "센서", Package: "robot/sensor"}, "`센서` 는 `robot/sensor` 에서 참조한다"},
	}
	for _, c := range cases {
		if got := renderSeeKo(c.s); got != c.want {
			t.Errorf("renderSeeKo(%+v) = %q, want %q", c.s, got, c.want)
		}
	}
}

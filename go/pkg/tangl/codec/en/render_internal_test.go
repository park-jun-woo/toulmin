//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderInternal — renderInternal covers on-event, every, and every-until heads
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderInternal checks the on-event, plain every, and every-until trigger heads.
func TestRenderInternal(t *testing.T) {
	cases := []struct {
		in   ast.Internal
		want string
	}{
		{ast.Internal{Kind: ast.OnEvent, Event: "start"}, "on start"},
		{ast.Internal{Kind: ast.EveryTick, Interval: "30s"}, "every 30s"},
		{ast.Internal{Kind: ast.EveryTick, Interval: "1s", Until: "can serve coffee"}, "every 1s until `can serve coffee`"},
	}
	for _, tc := range cases {
		got, err := renderInternal(tc.in)
		if err != nil {
			t.Errorf("renderInternal(%+v): %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("renderInternal = %q, want %q", got, tc.want)
		}
	}

	if _, err := renderInternal(ast.Internal{Kind: ast.InternalKind(99)}); err == nil {
		t.Error("renderInternal accepted an unknown InternalKind")
	}
}

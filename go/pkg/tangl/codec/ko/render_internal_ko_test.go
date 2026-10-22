//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderInternalKo — renderInternalKo covers 발생 시, 매 ~마다, and 매 ~마다 까지 heads
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderInternalKo checks the on-event head, a plain tick (with interval
// denormalization), and a tick with an until clause, plus the unknown-kind error.
func TestRenderInternalKo(t *testing.T) {
	cases := []struct {
		in   ast.Internal
		want string
	}{
		{ast.Internal{Kind: ast.OnEvent, Event: "버튼 눌림"}, "`버튼 눌림` 발생 시"},
		{ast.Internal{Kind: ast.EveryTick, Interval: "30s"}, "매 30초마다"},
		{ast.Internal{Kind: ast.EveryTick, Interval: "1s", Until: "서빙 가능"}, "매 1초마다, `서빙 가능` 까지"},
		{ast.Internal{Kind: ast.EveryTick, Interval: "7:00"}, "매 7:00마다"},
	}
	for _, c := range cases {
		got, err := renderInternalKo(c.in)
		if err != nil {
			t.Errorf("renderInternalKo(%+v): %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("renderInternalKo(%+v) = %q, want %q", c.in, got, c.want)
		}
	}
	if _, err := renderInternalKo(ast.Internal{Kind: ast.InternalKind(99)}); err == nil {
		t.Error("renderInternalKo accepted an unknown InternalKind")
	}
}

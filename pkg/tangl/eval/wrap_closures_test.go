//ff:func feature=tangl type=engine control=sequence
//ff:what TestWrapClosures — tests wrapping compiled closures as fresh-scope rule funcs
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWrapClosures(t *testing.T) {
	rules := []ast.InlineRule{{Name: "r"}}
	compiled := map[string]evalFunc{
		"r": func(sc *scope) (bool, any) {
			v, ok := sc.ctx.Get("x")
			return ok && len(sc.binders) == 0, v
		},
	}
	out := wrapClosures(rules, compiled)
	if len(out) != 1 {
		t.Fatalf("want 1 wrapped rule, got %d", len(out))
	}
	b, w := out["r"](testCtx(), nil)
	if !b || w != "1" {
		t.Fatalf("wrapped rule should run on a fresh seeded scope, got %v,%v", b, w)
	}
}

//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderClosureCall — tests holds-closure emission for bare and aliased calls
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderClosureCall(t *testing.T) {
	got := renderClosureCall(ast.Call{Fn: ast.Ref{Name: "is spam"}})
	if !strings.Contains(got, "b, ev := isSpam(ctx, specs)") || !strings.Contains(got, `tanglHolds(b, ev, "is spam")`) {
		t.Errorf("bare call closure: %q", got)
	}
	got = renderClosureCall(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}})
	if !strings.Contains(got, "b, ev := nlp.score(ctx, specs)") || !strings.Contains(got, `tanglHolds(b, ev, "nlp.score")`) {
		t.Errorf("aliased call closure: %q", got)
	}
}

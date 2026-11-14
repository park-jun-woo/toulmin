//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildRulesSkipsHostRules — host (uses) rules emit no Go function; expression rules do
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestBuildRulesSkipsHostRules(t *testing.T) {
	var w strings.Builder
	err := buildRules(&w, []ast.InlineRule{
		{Name: "host only", Using: &ast.Ref{Alias: "ratchet", Name: "triesExhausted"}},
		{Name: "big", Cond: ast.Compare{Field: "n", Op: "is greater than", Value: "3"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := w.String()
	if strings.Contains(out, "hostOnly") || strings.Contains(out, "triesExhausted") {
		t.Fatalf("host rule must not be rendered:\n%s", out)
	}
	if !strings.Contains(out, "func big(") {
		t.Fatalf("expression rule missing:\n%s", out)
	}
	w.Reset()
	if err := buildRules(&w, []ast.InlineRule{{Name: "h", Using: &ast.Ref{Name: "f"}}}); err != nil || w.Len() != 0 {
		t.Fatalf("only host rules: err=%v out=%q", err, w.String())
	}
}

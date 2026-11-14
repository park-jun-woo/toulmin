//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRuleHost — a host rule renders "uses" with optional with/grounded-on clauses and re-parses
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderRuleHost(t *testing.T) {
	bare := ast.InlineRule{Name: "tries exhausted", Using: &ast.Ref{Alias: "ratchet", Name: "triesExhausted"}}
	if got, want := renderRule(bare), "`tries exhausted` uses `ratchet`.`triesExhausted`"; got != want {
		t.Fatalf("bare host rule = %q, want %q", got, want)
	}
	full := ast.InlineRule{
		Name:    "tries exhausted",
		Using:   &ast.Ref{Alias: "ratchet", Name: "triesExhausted"},
		With:    []string{"max attempts", "window"},
		Grounds: []string{"attempts", "log"},
	}
	want := "`tries exhausted` uses `ratchet`.`triesExhausted` with `max attempts` and `window` grounded on `attempts` and `log`"
	if got := renderRule(full); got != want {
		t.Fatalf("full host rule = %q, want %q", got, want)
	}
}

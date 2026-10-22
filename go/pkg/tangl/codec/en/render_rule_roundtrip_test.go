//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRuleRoundtrip — round-trip nested Rules (either, block not, inline not) through the pipeline
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestRenderRuleRoundtrip drives the Rules multi-line renderer over an either
// group, a block "not" (a negated and-list), and an inline "not `f` op v" — the
// nesting shapes the fixtures do not all cover — and asserts encode -> render ->
// re-parse yields the same AST.
func TestRenderRuleRoundtrip(t *testing.T) {
	src := "## tangl:Subject\n- this document is `demo`\n\n" +
		"## tangl:Rules\n" +
		"1. `a` when\n" +
		"    - either\n" +
		"      - `p` is at least 1\n" +
		"      - and `q` equals \"x\"\n" +
		"    - or `r` is less than 2\n" +
		"- `b` when\n" +
		"    - not\n" +
		"      - `s` equals 1\n" +
		"      - and `t` is not empty\n" +
		"- `c` when\n" +
		"    - not `u` is empty\n\n" +
		"## tangl:Cases\n- in case of `x`\n  - `f` is required\n"

	d, err := ENCodec{}.Encode([]byte(src), "rule.md")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	rendered, err := doc.Render(d, ENCodec{})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	reparsed, err := parser.ParseSource(string(rendered), d.AST.Path)
	if err != nil {
		t.Fatalf("re-parse failed: %v\n---\n%s", err, rendered)
	}
	if !tangeul.EqualIgnoringLine(d.AST, reparsed) {
		t.Errorf("re-parsed AST differs\n--- rendered ---\n%s", rendered)
	}
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEither — a nested either/and/or Rules condition, not used by the spec's three examples
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseEither exercises the Rules section's nested either/and/or tree.
func TestParseEither(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Rules\n" +
		"1. `qualifies` when\n" +
		"    - either\n" +
		"      - `income` is at least 1000\n" +
		"      - and `score` is at least 700\n" +
		"    - or `vip` equals true\n\n" +
		"## tangl:Cases\n- in case of `x`\n  - `a` is required\n"
	doc, err := ParseSource(src, "either.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	if len(doc.Rules) != 1 {
		t.Fatalf("Rules = %+v, want 1", doc.Rules)
	}
	top, ok := doc.Rules[0].Cond.(ast.Logic)
	if !ok || top.Op != "or" {
		t.Fatalf("Cond = %#v, want top-level or Logic", doc.Rules[0].Cond)
	}
	either, ok := top.Terms[0].(ast.Either)
	if !ok || len(either.Terms) != 1 {
		t.Fatalf("Terms[0] = %#v, want Either", top.Terms[0])
	}
	inner, ok := either.Terms[0].(ast.Logic)
	if !ok || inner.Op != "and" {
		t.Fatalf("either inner = %#v, want and Logic", either.Terms[0])
	}
}

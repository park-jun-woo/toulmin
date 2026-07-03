//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitRule — encodes an inline rule's condition + Rule node, and propagates an EncodeExpr error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitRule covers the condition-encoding error path (an unrecognized
// Compare operator makes packet.EncodeExpr fail) alongside the success path.
// appendNode/appendLabel never fail here (NodeRule/PropLabel are always
// emitted with consistent flags), so their error-return lines are unreachable
// through any legitimate call and are not exercised.
func TestEmitRule(t *testing.T) {
	ref := doc.StmtRef{Sec: 4, Path: []int{0}}
	span := word.Span{}

	cases := []struct {
		name    string
		rule    ast.InlineRule
		wantErr bool
	}{
		{"valid condition", ast.InlineRule{Name: "ready", Cond: ast.Compare{Field: "level", Op: "equals", Value: "5"}}, false},
		{"unknown compare op", ast.InlineRule{Name: "bad", Cond: ast.Compare{Field: "level", Op: "bogus", Value: "5"}}, true},
	}

	check := func(t *testing.T, rule ast.InlineRule, wantErr bool) {
		e := &enc{tt: &tidTable{byRef: map[string]uint16{}}, nodeNext: 100}
		got, err := emitRule(e, ref, rule, span)
		if wantErr && err == nil {
			t.Fatalf("emitRule(%+v): expected error, got nil", rule)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("emitRule(%+v): %v", rule, err)
		}
		if len(got) == 0 {
			t.Errorf("emitRule(%+v): expected non-empty words", rule)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.rule, c.wantErr) })
	}
}

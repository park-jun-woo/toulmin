//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitRule — encodes an inline rule's condition + Rule node, and propagates an EncodeExpr error
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitRule covers the condition-encoding error path (an unrecognized
// Compare operator makes packet.EncodeExpr fail) alongside the success path.
// appendNode/appendLabel never fail here (NodeRule/PropLabel are always
// emitted with consistent flags), so their error-return lines are unreachable
// through any legitimate call and are not exercised.
func TestEmitRule(t *testing.T) {
	ref := doc.StmtRef{Sec: 4, Path: []int{0}}
	span := toulmin.Span{}

	cases := []struct {
		name    string
		rule    ast.InlineRule
		wantErr bool
	}{
		{"valid condition", ast.InlineRule{Name: "ready", Cond: ast.Compare{Field: "level", Op: "equals", Value: "5"}}, false},
		{"unknown compare op", ast.InlineRule{Name: "bad", Cond: ast.Compare{Field: "level", Op: "bogus", Value: "5"}}, true},
	}

	check := func(t *testing.T, rule ast.InlineRule, wantErr bool) {
		e := &enc{tt: &tidTable{byRef: map[string]toulmin.TID{}}, wire: testWire, next: 100}
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

	ruleRef := doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{0}}
	rule := ast.InlineRule{Name: "r", Cond: ast.Compare{Field: "x", Op: "is empty"}}
	if _, err := emitRule(testEnc(nodeTT(ruleRef, 0x10000)), ruleRef, rule, toulmin.Span{}); err == nil {
		t.Fatal("node id overflow: want error")
	}
	failsAtEachID(t, 3, nodeTT(ruleRef, 1), func(e *enc) error {
		_, err := emitRule(e, ruleRef, rule, toulmin.Span{})
		return err
	})
	host := ast.InlineRule{Name: "h", Using: &ast.Ref{Name: "f"}, Grounds: []string{"g"}}
	if ws, err := emitRule(testEnc(nodeTT(ruleRef, 1)), ruleRef, host, toulmin.Span{}); err != nil || len(ws) == 0 {
		t.Fatalf("host rule: %v", err)
	}
	failsAtEachID(t, 6, nodeTT(ruleRef, 1), func(e *enc) error {
		_, err := emitRule(e, ruleRef, host, toulmin.Span{})
		return err
	})
}

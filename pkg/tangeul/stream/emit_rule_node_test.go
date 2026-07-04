//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitRuleNode — role/using/checking/with/qualifier combinations and each fallible field's error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitRuleNode walks every optional field emitRuleNode may encode (Using,
// Checking, With, Qualified) both absent and present, plus the three fallible
// inputs: an unrecognized Role, an Using reference to an undeclared alias, a
// Checking reference to an undeclared case, and an out-of-range Qualified
// value. The final appendNode/appendLabel calls never fail here (their flags
// and matching data are always set together by emitRuleNode itself), so their
// error-return lines are unreachable through any legitimate call and are not
// exercised.
func TestEmitRuleNode(t *testing.T) {
	span := word.Span{}
	qualOK := 0.5
	qualBad := 2.0

	cases := []struct {
		name    string
		n       ast.Node
		wantErr bool
	}{
		{"minimal, no optional fields", ast.Node{Name: "n0", Role: ast.GeneralRule}, false},
		{"unknown role", ast.Node{Name: "n1", Role: ast.Role(99)}, true},
		{"using declared alias", ast.Node{Name: "n2", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "bank", Name: "hasBalance"}}, false},
		{"using local (no alias)", ast.Node{Name: "n2b", Role: ast.GeneralRule, Using: &ast.Ref{Name: "localFn"}}, false},
		{"using undeclared alias", ast.Node{Name: "n3", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "nope", Name: "f"}}, true},
		{"checking declared case", ast.Node{Name: "n4", Role: ast.GeneralRule, Checking: &ast.Ref{Name: "can check"}}, false},
		{"checking undeclared case", ast.Node{Name: "n5", Role: ast.GeneralRule, Checking: &ast.Ref{Name: "missing"}}, true},
		{"with terms", ast.Node{Name: "n6", Role: ast.GeneralRule, With: []string{"a", "b"}}, false},
		{"qualified ok", ast.Node{Name: "n7", Role: ast.GeneralRule, Qualified: &qualOK}, false},
		{"qualified out of range", ast.Node{Name: "n8", Role: ast.GeneralRule, Qualified: &qualBad}, true},
		{
			"all optional fields combined",
			ast.Node{
				Name:      "n9",
				Role:      ast.CounterRule,
				Using:     &ast.Ref{Alias: "bank", Name: "isFrozen"},
				Checking:  &ast.Ref{Name: "can check"},
				With:      []string{"term"},
				Qualified: &qualOK,
			},
			false,
		},
	}

	check := func(t *testing.T, n ast.Node, wantErr bool) {
		ref := doc.StmtRef{Sec: 5, Path: []int{0, 1, 0}}
		tt := &tidTable{
			byRef:       map[string]uint16{},
			aliasByName: map[string]uint16{"bank": 2},
			caseByName:  map[string]uint16{"can check": 9},
			nodeByCase:  map[string]uint16{},
		}
		e := &enc{tt: tt, nodeNext: 100}
		got, err := emitRuleNode(e, ref, n, span)
		if wantErr && err == nil {
			t.Fatalf("emitRuleNode(%+v): expected error, got nil", n)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("emitRuleNode(%+v): %v", n, err)
		}
		if len(got) == 0 {
			t.Errorf("emitRuleNode(%+v): expected non-empty words", n)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.n, c.wantErr) })
	}
}

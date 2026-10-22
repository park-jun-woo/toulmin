//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitCaseSlot — routes require/node/attack/exec slot kinds and errors on an unknown slotKind
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitCaseSlot drives every emitCaseSlot switch arm: slotKind 0 (Require),
// 1 (Node), 2 (Attack), 3 (Exec), 4 (Supersede), and an unrecognized slotKind
// that falls to the default hard error.
func TestEmitCaseSlot(t *testing.T) {
	e := testEnc(&tidTable{
		byRef: map[string]toulmin.TID{},
		nodeByCase: map[string]toulmin.TID{
			nodeKey(0, "atk-a"): 10,
			nodeKey(0, "atk-b"): 11,
			nodeKey(0, "trig"):  12,
		},
		caseByName: map[string]toulmin.TID{"target-case": 20},
	})
	span := toulmin.Span{Off: 0, Len: 1}

	cases := []struct {
		name    string
		path    []int
		elem    any
		wantErr bool
	}{
		{"require", []int{0, 0, 0}, ast.Require{Field: "amount"}, false},
		{"node", []int{0, 1, 0}, ast.Node{Name: "n", Role: ast.GeneralRule}, false},
		{"attack", []int{0, 2, 0}, ast.Attack{Attacker: "atk-a", Target: "atk-b"}, false},
		{"exec", []int{0, 3, 0}, ast.Exec{Kind: ast.RunExec, Node: "trig", Case: "target-case"}, false},
		{"supersede", []int{0, 4, 0}, ast.Supersede{Superseder: "atk-a", Superseded: "atk-b"}, false},
		{"unknown slotKind", []int{0, 5, 0}, nil, true},
	}
	for _, c := range cases {
		ref := doc.StmtRef{Sec: 5, Path: c.path}
		got, err := emitCaseSlot(e, ref, c.elem, span)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: emitCaseSlot: %v", c.name, err)
		}
		if len(got) == 0 {
			t.Errorf("%s: expected non-empty words", c.name)
		}
	}
}

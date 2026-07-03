//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitCaseSlot — routes require/node/attack/exec slot kinds and errors on an unknown slotKind
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitCaseSlot drives every emitCaseSlot switch arm: slotKind 0 (Require),
// 1 (Node), 2 (Attack), 3 (Exec), and an unrecognized slotKind that falls to the
// default hard error.
func TestEmitCaseSlot(t *testing.T) {
	e := &enc{tt: &tidTable{
		byRef: map[string]uint16{},
		nodeByCase: map[string]uint16{
			nodeKey(0, "atk-a"): 10,
			nodeKey(0, "atk-b"): 11,
			nodeKey(0, "trig"):  12,
		},
		caseByName: map[string]uint16{"target-case": 20},
	}}
	span := word.Span{TID: 0, Off: 0, Len: 1}

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
		{"unknown slotKind", []int{0, 4, 0}, nil, true},
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

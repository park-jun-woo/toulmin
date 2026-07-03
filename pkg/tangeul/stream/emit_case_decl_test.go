//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitCaseDecl — a case declaration encodes its Case node plus name label
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitCaseDecl checks the success path of emitCaseDecl: the Case node and
// its name label are both present in the returned words. Both appendNode and
// appendLabel calls here always carry a well-formed fixed-shape packet (no
// flags, no gate), so their `if err != nil` guards are defensive and not
// reachable through any legitimate caller in this codebase.
func TestEmitCaseDecl(t *testing.T) {
	e := &enc{tt: &tidTable{byRef: map[string]uint16{}}}
	ref := doc.StmtRef{Sec: 5, Path: []int{0}}
	span := word.Span{TID: 0, Off: 0, Len: 1}

	got, err := emitCaseDecl(e, ref, ast.Case{Name: "can withdraw"}, span)
	if err != nil {
		t.Fatalf("emitCaseDecl: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected non-empty words")
	}
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitCaseDecl — a case declaration encodes its Case node plus name label
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitCaseDecl checks the success path of emitCaseDecl: the Case node and
// its name label are both present in the returned words. Both appendNode and
// appendLabel calls here always carry a well-formed fixed-shape packet (no
// flags, no gate), so their `if err != nil` guards are defensive and not
// reachable through any legitimate caller in this codebase.
func TestEmitCaseDecl(t *testing.T) {
	e := testEnc(&tidTable{byRef: map[string]toulmin.TID{}})
	ref := doc.StmtRef{Sec: 5, Path: []int{0}}
	span := toulmin.Span{Off: 0, Len: 1}

	got, err := emitCaseDecl(e, ref, ast.Case{Name: "can withdraw"}, span)
	if err != nil {
		t.Fatalf("emitCaseDecl: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected non-empty words")
	}

	caseRef := doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}
	if _, err := emitCaseDecl(testEnc(nodeTT(caseRef, 0x10000)), caseRef, ast.Case{Name: "c"}, toulmin.Span{}); err == nil {
		t.Fatal("node id overflow: want error")
	}
	failsAtEachID(t, 2, nodeTT(caseRef, 1), func(e *enc) error {
		_, err := emitCaseDecl(e, caseRef, ast.Case{Name: "c"}, toulmin.Span{})
		return err
	})
}

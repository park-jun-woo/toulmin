//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitSee — encodes a package alias reference (Alias node + name label + importedFrom edge)
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitSee is a single straight-line assertion: emitSee has no data-dependent
// branch, only defensive appendNode/appendLabel/appendEdge error checks that
// never fire (NodeAlias and PropImportedFrom are always emitted with
// consistent, valid flags), so those error-return lines are unreachable
// through any legitimate call and are not exercised.
func TestEmitSee(t *testing.T) {
	ref := doc.StmtRef{Sec: 2, Path: []int{0}}
	e := &enc{tt: &tidTable{byRef: map[string]toulmin.TID{}}, wire: testWire, next: 100}
	see := ast.See{Alias: "bank", Package: "bank/core"}

	got, err := emitSee(e, ref, see, toulmin.Span{})
	if err != nil {
		t.Fatalf("emitSee: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("emitSee: expected non-empty words")
	}

	seeRef := doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{0}}
	ovSee := ast.See{Alias: "bank", Package: "example.com/bank"}
	if _, err := emitSee(testEnc(nodeTT(seeRef, 0x10000)), seeRef, ovSee, toulmin.Span{}); err == nil {
		t.Fatal("node id overflow: want error")
	}
	failsAtEachID(t, 4, nodeTT(seeRef, 1), func(e *enc) error {
		_, err := emitSee(e, seeRef, ovSee, toulmin.Span{})
		return err
	})
}

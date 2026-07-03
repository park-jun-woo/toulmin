//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitSee — encodes a package alias reference (Alias node + name label + importedFrom edge)
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitSee is a single straight-line assertion: emitSee has no data-dependent
// branch, only defensive appendNode/appendLabel/appendEdge error checks that
// never fire (NodeAlias and PropImportedFrom are always emitted with
// consistent, valid flags), so those error-return lines are unreachable
// through any legitimate call and are not exercised.
func TestEmitSee(t *testing.T) {
	ref := doc.StmtRef{Sec: 2, Path: []int{0}}
	e := &enc{tt: &tidTable{byRef: map[string]uint16{}}, nodeNext: 100}
	see := ast.See{Alias: "bank", Package: "bank/core"}

	got, err := emitSee(e, ref, see, word.Span{})
	if err != nil {
		t.Fatalf("emitSee: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("emitSee: expected non-empty words")
	}
}

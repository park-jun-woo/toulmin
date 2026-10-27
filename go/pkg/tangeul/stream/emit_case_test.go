//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitCase — a one-element path routes to the case decl, a three-element path routes to a body slot
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitCase checks both emitCase branches: len(Path)==1 dispatches to
// emitCaseDecl with the elem asserted as ast.Case, and any other Path length
// (a three-element case-body coordinate) dispatches to emitCaseSlot.
func TestEmitCase(t *testing.T) {
	e := testEnc(&tidTable{})
	span := toulmin.Span{Off: 0, Len: 1}

	t.Run("decl", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 5, Path: []int{0}}
		got, err := emitCase(e, ref, ast.Case{Name: "can withdraw"}, span)
		if err != nil {
			t.Fatalf("emitCase decl: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})

	t.Run("slot", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 5, Path: []int{0, 1, 0}}
		got, err := emitCase(e, ref, ast.Node{Name: "n", Role: ast.GeneralRule}, span)
		if err != nil {
			t.Fatalf("emitCase slot: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})
}

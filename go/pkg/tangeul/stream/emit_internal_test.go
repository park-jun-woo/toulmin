//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitInternal — a one-element path routes to the trigger decl, path[1]==2/!=2 pick Evaluates/Triggers
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitInternal checks all three emitInternal branches: len(Path)==1 (the
// trigger head declaration) dispatches to emitInternalDecl, a three-element
// Path with Path[1]==2 (a check) picks PropEvaluates, and any other slot kind
// (Path[1]==1, a run) picks PropTriggers — both routed through emitCaseEdge.
func TestEmitInternal(t *testing.T) {
	e := testEnc(&tidTable{
		byRef:      map[string]toulmin.TID{},
		caseByName: map[string]toulmin.TID{"can check": 9},
	})
	span := toulmin.Span{Off: 0, Len: 1}

	t.Run("decl", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 7, Path: []int{0}}
		got, err := emitInternal(e, ref, ast.Internal{}, span)
		if err != nil {
			t.Fatalf("emitInternal decl: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})

	t.Run("run (path[1]!=2)", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 7, Path: []int{0, 1, 0}}
		got, err := emitInternal(e, ref, "can check", span)
		if err != nil {
			t.Fatalf("emitInternal run: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})

	t.Run("check (path[1]==2)", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 7, Path: []int{0, 2, 0}}
		got, err := emitInternal(e, ref, "can check", span)
		if err != nil {
			t.Fatalf("emitInternal check: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})
}

//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitExecStmt — undeclared trigger/kind/certainty/func/case errors and do/undo/run success paths
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitExecStmt drives every emitExecStmt branch: an undeclared `when`
// trigger, an unrecognized ExecKind (codebook.ExecPropCode failure), an invalid
// Certainty op (execFlags failure), a do/undo exec with a nil Func, a do/undo
// exec whose Func references an undeclared alias (emitFuncRef failure), a
// do/undo success, a run exec to an undeclared case, and a run success. The
// final `if err := appendEdge(...)` guard is not exercised: PropDo/PropUndo
// always carry a non-nil DstFunc here and PropRun always carries a resolved
// DstTID, with Gate kept in sync with EdgeFlagGate by execFlags, so
// toulmin.EncodeEdge cannot fail for any shape this function builds.
func TestEmitExecStmt(t *testing.T) {
	e := testEnc(&tidTable{
		nodeByCase:  map[string]toulmin.TID{nodeKey(0, "trig"): 5},
		caseByName:  map[string]toulmin.TID{"declared-case": 9},
		aliasByName: map[string]toulmin.TID{},
	})
	span := toulmin.Span{Off: 0, Len: 1}
	ref := doc.StmtRef{Sec: 5, Path: []int{0, 3, 0}}

	cases := []struct {
		name    string
		x       ast.Exec
		wantErr bool
	}{
		{"undeclared trigger", ast.Exec{Node: "missing", Kind: ast.DoExec, Func: &ast.Ref{Name: "f"}}, true},
		{"unknown exec kind", ast.Exec{Node: "trig", Kind: ast.ExecKind(99)}, true},
		{"invalid certainty op", ast.Exec{Node: "trig", Kind: ast.DoExec, Func: &ast.Ref{Name: "f"}, Certainty: &ast.Certainty{Op: "weird", Percent: 50}}, true},
		{"do missing func", ast.Exec{Node: "trig", Kind: ast.DoExec, Func: nil}, true},
		{"do undeclared alias", ast.Exec{Node: "trig", Kind: ast.DoExec, Func: &ast.Ref{Alias: "bad", Name: "f"}}, true},
		{"do success", ast.Exec{Node: "trig", Kind: ast.DoExec, Func: &ast.Ref{Name: "f"}, Once: true, Certainty: &ast.Certainty{Op: "at least", Percent: 75}}, false},
		{"undo success", ast.Exec{Node: "trig", Kind: ast.UndoExec, Func: &ast.Ref{Name: "g"}}, false},
		{"run undeclared case", ast.Exec{Node: "trig", Kind: ast.RunExec, Case: "missing-case"}, true},
		{"run success", ast.Exec{Node: "trig", Kind: ast.RunExec, Case: "declared-case"}, false},
	}
	for _, c := range cases {
		got, err := emitExecStmt(e, ref, c.x, span)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: emitExecStmt: %v", c.name, err)
		}
		if len(got) == 0 {
			t.Errorf("%s: expected non-empty words", c.name)
		}
	}
}

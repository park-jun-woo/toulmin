//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteNodeExecs — tests writeNodeExecs for empty, RunExec-skip, and do/undo/do interleaving order
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWriteNodeExecs(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		writeNodeExecs(&w, "subj", "case1", "node1", nil)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("skips RunExec", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{{Kind: ast.RunExec, Case: "other", Node: "node1"}}
		writeNodeExecs(&w, "subj", "case1", "node1", execs)
		if w.String() != "" {
			t.Errorf("expected empty output for RunExec entries, got %q", w.String())
		}
	})

	t.Run("arms undo between two dos in document order", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Kind: ast.DoExec, Once: true, Func: &ast.Ref{Name: "actA"}, Node: "node1"},
			{Kind: ast.UndoExec, Func: &ast.Ref{Name: "undoA"}, Node: "node1"},
			{Kind: ast.DoExec, Func: &ast.Ref{Name: "actB"}, Node: "node1"},
		}
		writeNodeExecs(&w, "subj", "case1", "node1", execs)
		got := w.String()

		idxA := strings.Index(got, "actA(t.Ctx())")
		idxUndo := strings.Index(got, "return undoA(t.Ctx())")
		idxB := strings.Index(got, "actB(t.Ctx())")
		if idxA < 0 || idxUndo < 0 || idxB < 0 {
			t.Fatalf("expected actA, undoA push, and actB all present, got %q", got)
		}
		if !(idxA < idxUndo && idxUndo < idxB) {
			t.Errorf("expected order actA < undoA push < actB, got positions %d, %d, %d in %q", idxA, idxUndo, idxB, got)
		}
		if !strings.Contains(got, `"once:subj.case1.node1#0"`) {
			t.Errorf("expected once key indexed by do-position 0, got %q", got)
		}
	})
}

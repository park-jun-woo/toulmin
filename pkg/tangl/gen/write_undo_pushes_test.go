//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteUndoPushes — tests writeUndoPushes for skip-non-UndoExec and matching-UndoExec branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWriteUndoPushes(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		writeUndoPushes(&w, nil)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("skips non-UndoExec", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Kind: ast.DoExec, Func: &ast.Ref{Name: "act1"}, Node: "node1"},
			{Kind: ast.RunExec, Case: "other", Node: "node1"},
		}
		writeUndoPushes(&w, execs)
		if w.String() != "" {
			t.Errorf("expected empty output for non-UndoExec entries, got %q", w.String())
		}
	})

	t.Run("with UndoExec", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Kind: ast.DoExec, Func: &ast.Ref{Name: "act1"}, Node: "node1"},
			{Kind: ast.UndoExec, Func: &ast.Ref{Name: "cleanup1"}, Node: "node1"},
			{Kind: ast.UndoExec, Func: &ast.Ref{Name: "cleanup2"}, Node: "node1"},
		}
		writeUndoPushes(&w, execs)
		got := w.String()
		if !strings.Contains(got, "tangl.PushCompensation(t.Ctx(), func() error {") {
			t.Errorf("expected push compensation call, got %q", got)
		}
		if !strings.Contains(got, "return cleanup1(t.Ctx())") {
			t.Errorf("expected cleanup1 call, got %q", got)
		}
		if !strings.Contains(got, "return cleanup2(t.Ctx())") {
			t.Errorf("expected cleanup2 call, got %q", got)
		}
	})
}

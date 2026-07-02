//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteDoStatements — tests writeDoStatements for skip-non-DoExec and matching-DoExec branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWriteDoStatements(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		writeDoStatements(&w, "subj", "case1", "node1", nil)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("skips non-DoExec", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Kind: ast.UndoExec, Func: &ast.Ref{Name: "cleanup"}, Node: "node1"},
			{Kind: ast.RunExec, Case: "other", Node: "node1"},
		}
		writeDoStatements(&w, "subj", "case1", "node1", execs)
		if w.String() != "" {
			t.Errorf("expected empty output for non-DoExec entries, got %q", w.String())
		}
	})

	t.Run("with DoExec", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Kind: ast.UndoExec, Func: &ast.Ref{Name: "cleanup"}, Node: "node1"},
			{Kind: ast.DoExec, Func: &ast.Ref{Name: "act1"}, Node: "node1"},
			{Kind: ast.DoExec, Func: &ast.Ref{Name: "act2"}, Node: "node1"},
		}
		writeDoStatements(&w, "subj", "case1", "node1", execs)
		got := w.String()
		if !strings.Contains(got, "act1(t.Ctx())") {
			t.Errorf("expected act1 call, got %q", got)
		}
		if !strings.Contains(got, "act2(t.Ctx())") {
			t.Errorf("expected act2 call, got %q", got)
		}
	})
}

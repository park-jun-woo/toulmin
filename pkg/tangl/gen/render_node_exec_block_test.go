//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderNodeExecBlock — tests renderNodeExecBlock for do/undo-block and run-attachment branch combinations plus error propagation
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderNodeExecBlock(t *testing.T) {
	ni := nodeInfo{Var: "nA", Node: ast.Node{Name: "nodeA"}}

	t.Run("no do/undo no run", func(t *testing.T) {
		var w strings.Builder
		err := renderNodeExecBlock(&w, "subj", "caseA", ni, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.String() != "\tnA\n" {
			t.Errorf("got %q, want %q", w.String(), "\tnA\n")
		}
	})

	t.Run("do/undo without run", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{{Node: "nodeA", Kind: ast.DoExec, Func: &ast.Ref{Name: "doIt"}}}
		err := renderNodeExecBlock(&w, "subj", "caseA", ni, execs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if !strings.Contains(out, ".RunOn(func(self toulmin.TraceEntry, t toulmin.Trace) error {") {
			t.Errorf("expected RunOn block: %q", out)
		}
		if strings.Contains(out, ".Run(") {
			t.Errorf("unexpected Run attachment: %q", out)
		}
	})

	t.Run("run without do/undo", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{{Node: "nodeA", Kind: ast.RunExec, Case: "otherCase"}}
		err := renderNodeExecBlock(&w, "subj", "caseA", ni, execs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if strings.Contains(out, ".RunOn(") {
			t.Errorf("unexpected RunOn block: %q", out)
		}
		if !strings.Contains(out, ".Run(otherCaseGraph)") {
			t.Errorf("expected Run attachment: %q", out)
		}
	})

	t.Run("do/undo and run together", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Node: "nodeA", Kind: ast.DoExec, Func: &ast.Ref{Name: "doIt"}},
			{Node: "nodeA", Kind: ast.RunExec, Case: "otherCase"},
		}
		err := renderNodeExecBlock(&w, "subj", "caseA", ni, execs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if !strings.Contains(out, ".RunOn(func(self toulmin.TraceEntry, t toulmin.Trace) error {") {
			t.Errorf("expected RunOn block: %q", out)
		}
		if !strings.Contains(out, ".Run(otherCaseGraph)") {
			t.Errorf("expected Run attachment: %q", out)
		}
	})

	t.Run("error from runTarget propagates", func(t *testing.T) {
		var w strings.Builder
		execs := []ast.Exec{
			{Node: "nodeA", Kind: ast.RunExec, Case: "caseOne"},
			{Node: "nodeA", Kind: ast.RunExec, Case: "caseTwo"},
		}
		err := renderNodeExecBlock(&w, "subj", "caseA", ni, execs)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

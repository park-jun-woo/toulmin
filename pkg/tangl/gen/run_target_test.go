//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRunTarget — tests runTarget for no-run-exec, single-run-exec, non-run-exec-skip, and duplicate-run-exec error branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRunTarget(t *testing.T) {
	t.Run("no execs returns empty target", func(t *testing.T) {
		got, err := runTarget("caseA", "nodeA", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "" {
			t.Errorf("got %q, want empty", got)
		}
	})

	t.Run("non-run exec is skipped", func(t *testing.T) {
		execs := []ast.Exec{{Kind: ast.DoExec}}
		got, err := runTarget("caseA", "nodeA", execs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "" {
			t.Errorf("got %q, want empty", got)
		}
	})

	t.Run("single run exec returns its case", func(t *testing.T) {
		execs := []ast.Exec{{Kind: ast.DoExec}, {Kind: ast.RunExec, Case: "targetCase"}}
		got, err := runTarget("caseA", "nodeA", execs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "targetCase" {
			t.Errorf("got %q, want %q", got, "targetCase")
		}
	})

	t.Run("more than one run exec is an error", func(t *testing.T) {
		execs := []ast.Exec{
			{Kind: ast.RunExec, Case: "caseOne"},
			{Kind: ast.RunExec, Case: "caseTwo"},
		}
		_, err := runTarget("caseA", "nodeA", execs)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

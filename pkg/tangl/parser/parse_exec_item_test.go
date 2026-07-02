//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseExecItem — tests parseExecItem for do/do-not/undo/run/unmatched branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseExecItem(t *testing.T) {
	t.Run("DoNot", func(t *testing.T) {
		it := item{Text: "do not `fn` when `a`", Line: 1}
		exec, matched, err := parseExecItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for 'do not'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec != (ast.Exec{}) {
			t.Errorf("expected zero-value Exec, got %+v", exec)
		}
	})

	t.Run("Do", func(t *testing.T) {
		it := item{Text: "do `fn` when `a`", Line: 2}
		exec, matched, err := parseExecItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true for 'do'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Kind != ast.DoExec {
			t.Errorf("expected DoExec, got %v", exec.Kind)
		}
	})

	t.Run("Undo", func(t *testing.T) {
		it := item{Text: "undo `fn` when `a`", Line: 3}
		exec, matched, err := parseExecItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true for 'undo'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Kind != ast.UndoExec {
			t.Errorf("expected UndoExec, got %v", exec.Kind)
		}
	})

	t.Run("Run", func(t *testing.T) {
		it := item{Text: "run `case1` when `a`", Line: 4}
		exec, matched, err := parseExecItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true for 'run'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Kind != ast.RunExec {
			t.Errorf("expected RunExec, got %v", exec.Kind)
		}
	})

	t.Run("Unmatched", func(t *testing.T) {
		it := item{Text: "something else entirely", Line: 5}
		exec, matched, err := parseExecItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for an unrecognized statement")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec != (ast.Exec{}) {
			t.Errorf("expected zero-value Exec, got %+v", exec)
		}
	})
}

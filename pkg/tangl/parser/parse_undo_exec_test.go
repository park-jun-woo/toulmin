//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUndoExec — tests parseUndoExec for ref/when/node/trailing/success branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseUndoExec(t *testing.T) {
	t.Run("BadRef", func(t *testing.T) {
		_, matched, err := parseUndoExec("not-a-ref when `a`", 1, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing function reference")
		}
		if !strings.Contains(err.Error(), "expected function reference after 'undo'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoWhen", func(t *testing.T) {
		_, matched, err := parseUndoExec("`fn` triggered", 2, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing 'when'")
		}
		if !strings.Contains(err.Error(), "expected 'when' in undo edge") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoNode", func(t *testing.T) {
		_, matched, err := parseUndoExec("`fn` when notbacktick", 3, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted node name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted node name after 'when'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, matched, err := parseUndoExec("`fn` when `a` extra", 4, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		exec, matched, err := parseUndoExec("`fn` when `a`", 5, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Kind != ast.UndoExec {
			t.Errorf("expected UndoExec, got %v", exec.Kind)
		}
		if exec.Node != "a" {
			t.Errorf("expected Node=a, got %q", exec.Node)
		}
		if exec.Func == nil || exec.Func.Name != "fn" {
			t.Errorf("expected Func.Name=fn, got %+v", exec.Func)
		}
	})
}

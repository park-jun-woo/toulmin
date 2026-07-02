//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDoExec — tests parseDoExec for ref/once/when/node/certainty branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseDoExec(t *testing.T) {
	t.Run("BadRef", func(t *testing.T) {
		_, matched, err := parseDoExec("not-a-ref when `a`", 1, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing function reference")
		}
		if !strings.Contains(err.Error(), "expected function reference after 'do'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoWhen", func(t *testing.T) {
		_, matched, err := parseDoExec("`fn` triggered", 2, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing 'when'")
		}
		if !strings.Contains(err.Error(), "expected 'when' in do edge") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoNode", func(t *testing.T) {
		_, matched, err := parseDoExec("`fn` when notbacktick", 3, "test.md")
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

	t.Run("NoCertaintyNoOnce", func(t *testing.T) {
		exec, matched, err := parseDoExec("`fn` when `a`", 4, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Kind != ast.DoExec {
			t.Errorf("expected DoExec, got %v", exec.Kind)
		}
		if exec.Once {
			t.Error("expected Once=false")
		}
		if exec.Node != "a" {
			t.Errorf("expected Node=a, got %q", exec.Node)
		}
		if exec.Certainty != nil {
			t.Errorf("expected nil Certainty, got %+v", exec.Certainty)
		}
		if exec.Func == nil || exec.Func.Name != "fn" {
			t.Errorf("expected Func.Name=fn, got %+v", exec.Func)
		}
	})

	t.Run("OnceCertaintyError", func(t *testing.T) {
		_, matched, err := parseDoExec("`fn` once when `a` bogus clause", 5, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error from parseCertaintyClause")
		}
		if !strings.Contains(err.Error(), "expected 'if' certainty clause") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("OnceCertaintySuccess", func(t *testing.T) {
		exec, matched, err := parseDoExec("`fn` once when `a` if at least 75% certain", 6, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !exec.Once {
			t.Error("expected Once=true")
		}
		if exec.Node != "a" {
			t.Errorf("expected Node=a, got %q", exec.Node)
		}
		if exec.Certainty == nil {
			t.Fatal("expected non-nil Certainty")
		}
		if exec.Certainty.Op != "at least" || exec.Certainty.Percent != 75 {
			t.Errorf("unexpected certainty: %+v", exec.Certainty)
		}
	})
}

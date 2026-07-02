//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRunExec — tests parseRunExec for case/when/node/trailing/success branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseRunExec(t *testing.T) {
	t.Run("NoCase", func(t *testing.T) {
		_, matched, err := parseRunExec("notbacktick when `a`", 1, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted case name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted case name after 'run'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoWhen", func(t *testing.T) {
		_, matched, err := parseRunExec("`case1` triggered", 2, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing 'when'")
		}
		if !strings.Contains(err.Error(), "expected 'when' in run edge") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoNode", func(t *testing.T) {
		_, matched, err := parseRunExec("`case1` when notbacktick", 3, "test.md")
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
		_, matched, err := parseRunExec("`case1` when `a` extra", 4, "test.md")
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
		exec, matched, err := parseRunExec("`case1` when `a`", 5, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Kind != ast.RunExec {
			t.Errorf("expected RunExec, got %v", exec.Kind)
		}
		if exec.Case != "case1" || exec.Node != "a" {
			t.Errorf("expected Case=case1 Node=a, got %+v", exec)
		}
	})
}

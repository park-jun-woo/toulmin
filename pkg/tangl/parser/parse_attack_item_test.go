//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseAttackItem — tests parseAttackItem for don't/do-not prefix dispatch, backtick/keyword error branches, and success
package parser

import (
	"strings"
	"testing"
)

func TestParseAttackItem(t *testing.T) {
	t.Run("NotAttack", func(t *testing.T) {
		_, ok, err := parseAttackItem(item{Text: "something else", Line: 1}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for a non-attack statement")
		}
	})

	t.Run("DoWithoutNot", func(t *testing.T) {
		_, ok, err := parseAttackItem(item{Text: "do something", Line: 2}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for 'do' without 'not'")
		}
	})

	t.Run("MissingTarget", func(t *testing.T) {
		_, ok, err := parseAttackItem(item{Text: "don't notabacktick", Line: 3}, "test.md")
		if !ok {
			t.Fatal("expected ok=true once don't prefix matched")
		}
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted target after don't/do not") {
			t.Fatalf("expected target error, got %v", err)
		}
	})

	t.Run("MissingWhen", func(t *testing.T) {
		_, ok, err := parseAttackItem(item{Text: "don't `target`", Line: 4}, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err == nil || !strings.Contains(err.Error(), "expected 'when' in don't edge") {
			t.Fatalf("expected when error, got %v", err)
		}
	})

	t.Run("MissingAttacker", func(t *testing.T) {
		_, ok, err := parseAttackItem(item{Text: "don't `target` when notabacktick", Line: 5}, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted attacker after 'when'") {
			t.Fatalf("expected attacker error, got %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, ok, err := parseAttackItem(item{Text: "don't `target` when `attacker` extra", Line: 6}, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Fatalf("expected trailing text error, got %v", err)
		}
	})

	t.Run("DontSuccess", func(t *testing.T) {
		atk, ok, err := parseAttackItem(item{Text: "don't `target` when `attacker`", Line: 7}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if atk.Target != "target" || atk.Attacker != "attacker" {
			t.Fatalf("expected Target/Attacker set, got %+v", atk)
		}
	})

	t.Run("DoNotSuccess", func(t *testing.T) {
		atk, ok, err := parseAttackItem(item{Text: "do not `target` when `attacker`", Line: 8}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if atk.Target != "target" || atk.Attacker != "attacker" {
			t.Fatalf("expected Target/Attacker set, got %+v", atk)
		}
	})
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestTakeBacktick — tests takeBacktick for no-prefix, unterminated, and success branches
package parser

import "testing"

func TestTakeBacktick(t *testing.T) {
	t.Run("NoPrefix", func(t *testing.T) {
		name, rest, ok := takeBacktick("  notbacktick")
		if ok {
			t.Fatal("expected ok=false for text with no leading backtick")
		}
		if name != "" {
			t.Errorf("expected empty name, got %q", name)
		}
		if rest != "notbacktick" {
			t.Errorf("expected rest='notbacktick', got %q", rest)
		}
	})

	t.Run("Unterminated", func(t *testing.T) {
		name, rest, ok := takeBacktick("`unterminated")
		if ok {
			t.Fatal("expected ok=false for an unterminated backtick")
		}
		if name != "" {
			t.Errorf("expected empty name, got %q", name)
		}
		if rest != "`unterminated" {
			t.Errorf("expected rest='`unterminated', got %q", rest)
		}
	})

	t.Run("Success", func(t *testing.T) {
		name, rest, ok := takeBacktick("  `credit` extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if name != "credit" {
			t.Errorf("expected name='credit', got %q", name)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})
}

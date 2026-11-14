//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRunCheckName — tests parseRunCheckName for missing-backtick, trailing-text, and success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseRunCheckName(t *testing.T) {
	t.Run("NoBacktick", func(t *testing.T) {
		_, _, ok, err := parseRunCheckName(" notbacktick", "run", 4, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted case name")
		}
		if !strings.Contains(err.Error(), `expected backtick-quoted case name after "run"`) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, _, ok, err := parseRunCheckName(" `case1` extra", "check", 5, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		name, kind, ok, err := parseRunCheckName(" `case1`", "run", 6, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "case1" || kind != "run" {
			t.Errorf("expected name=case1 kind=run, got name=%q kind=%q", name, kind)
		}
	})
}

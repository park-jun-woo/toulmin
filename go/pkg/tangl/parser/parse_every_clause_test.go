//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEveryClause — tests parseEveryClause for no-until and until-clause branches
package parser

import (
	"strings"
	"testing"
)

func TestParseEveryClause(t *testing.T) {
	t.Run("NoUntilEmptyInterval", func(t *testing.T) {
		_, _, err := parseEveryClause("   ", "test.md", 1)
		if err == nil {
			t.Fatal("expected an error for an empty interval")
		}
		if !strings.Contains(err.Error(), "expected interval after 'every'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoUntilSuccess", func(t *testing.T) {
		interval, caseName, err := parseEveryClause("  1 day  ", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if interval != "1 day" {
			t.Errorf("expected interval='1 day', got %q", interval)
		}
		if caseName != "" {
			t.Errorf("expected empty caseName, got %q", caseName)
		}
	})

	t.Run("UntilEmptyInterval", func(t *testing.T) {
		_, _, err := parseEveryClause("  until `done`", "test.md", 3)
		if err == nil {
			t.Fatal("expected an error for an empty interval before 'until'")
		}
		if !strings.Contains(err.Error(), "expected interval after 'every'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("UntilNoBacktick", func(t *testing.T) {
		_, _, err := parseEveryClause("1 day until notbacktick", "test.md", 4)
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted case name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted case name after 'until'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("UntilTrailingText", func(t *testing.T) {
		_, _, err := parseEveryClause("1 day until `done` extra", "test.md", 5)
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("UntilSuccess", func(t *testing.T) {
		interval, caseName, err := parseEveryClause("1 day until `done`", "test.md", 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if interval != "1 day" {
			t.Errorf("expected interval='1 day', got %q", interval)
		}
		if caseName != "done" {
			t.Errorf("expected caseName='done', got %q", caseName)
		}
	})
}

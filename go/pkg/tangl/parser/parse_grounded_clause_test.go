//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseGroundedClause — tests parseGroundedClause for missing-on, single-name, and-chain, and missing-backtick branches
package parser

import (
	"strings"
	"testing"
)

func TestParseGroundedClause(t *testing.T) {
	t.Run("MissingOn", func(t *testing.T) {
		_, _, err := parseGroundedClause("`g1`", "test.md", 1)
		if err == nil {
			t.Fatal("expected an error for a missing 'on' keyword")
		}
		if !strings.Contains(err.Error(), "expected 'on' after 'grounded'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("SingleName", func(t *testing.T) {
		names, rest, err := parseGroundedClause("on `source-body` extra", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 1 || names[0] != "source-body" {
			t.Errorf("expected names=[source-body], got %v", names)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("AndChain", func(t *testing.T) {
		names, rest, err := parseGroundedClause("on `source-body` and `mx` and `spf` extra", "test.md", 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 3 || names[0] != "source-body" || names[1] != "mx" || names[2] != "spf" {
			t.Errorf("expected names=[source-body mx spf], got %v", names)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("MissingBacktick", func(t *testing.T) {
		_, _, err := parseGroundedClause("on notbacktick", "test.md", 4)
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted name after 'grounded on'/'and'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("MidChainError", func(t *testing.T) {
		_, _, err := parseGroundedClause("on `source-body` and notbacktick", "test.md", 5)
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted name after 'and'")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted name after 'grounded on'/'and'") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

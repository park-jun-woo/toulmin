//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseWithClause — tests parseWithClause for the initial-error, single-term, and-chain, and mid-chain-error branches
package parser

import (
	"strings"
	"testing"
)

func TestParseWithClause(t *testing.T) {
	t.Run("InitialError", func(t *testing.T) {
		_, _, err := parseWithClause("notbacktick", "test.md", 1)
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted term")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted term after 'with'/'and'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("SingleTerm", func(t *testing.T) {
		terms, rest, err := parseWithClause("`term1` extra", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(terms) != 1 || terms[0] != "term1" {
			t.Errorf("expected terms=[term1], got %v", terms)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("AndChain", func(t *testing.T) {
		terms, rest, err := parseWithClause("`term1` and `term2` and `term3` extra", "test.md", 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(terms) != 3 || terms[0] != "term1" || terms[1] != "term2" || terms[2] != "term3" {
			t.Errorf("expected terms=[term1 term2 term3], got %v", terms)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("MidChainError", func(t *testing.T) {
		_, _, err := parseWithClause("`term1` and notbacktick", "test.md", 4)
		if err == nil {
			t.Fatal("expected an error for a missing backtick-quoted term after 'and'")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted term after 'with'/'and'") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUsingClause — tests parseUsingClause for ref/with branches
package parser

import (
	"strings"
	"testing"
)

func TestParseUsingClause(t *testing.T) {
	t.Run("BadRef", func(t *testing.T) {
		_, _, _, _, err := parseUsingClause("not-a-ref", "test.md", 1)
		if err == nil {
			t.Fatal("expected an error for a missing function reference")
		}
		if !strings.Contains(err.Error(), "expected function reference after 'using'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Bare", func(t *testing.T) {
		ref, with, qual, rest, err := parseUsingClause("`fn`", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if with != nil {
			t.Errorf("expected nil with, got %v", with)
		}
		if qual != nil {
			t.Errorf("expected nil qual, got %v", qual)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})

	t.Run("WithError", func(t *testing.T) {
		_, _, _, _, err := parseUsingClause("`fn` with notbacktick", "test.md", 3)
		if err == nil {
			t.Fatal("expected an error from parseWithClause")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted term after 'with'/'and'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("WithSuccess", func(t *testing.T) {
		ref, with, qual, rest, err := parseUsingClause("`fn` with `term1`", "test.md", 4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if len(with) != 1 || with[0] != "term1" {
			t.Errorf("expected with=[term1], got %v", with)
		}
		if qual != nil {
			t.Errorf("expected nil qual, got %v", qual)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})
}

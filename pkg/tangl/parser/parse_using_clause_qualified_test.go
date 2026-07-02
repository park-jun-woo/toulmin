//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUsingClause_Qualified — tests parseUsingClause for qualified/with-and-qualified branches
package parser

import (
	"strings"
	"testing"
)

func TestParseUsingClause_Qualified(t *testing.T) {
	t.Run("QualifiedError", func(t *testing.T) {
		_, _, _, _, err := parseUsingClause("`fn` qualified abc", "test.md", 5)
		if err == nil {
			t.Fatal("expected an error from parseQualifiedClause")
		}
		if !strings.Contains(err.Error(), "expected float after 'qualified'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("QualifiedSuccess", func(t *testing.T) {
		ref, with, qual, rest, err := parseUsingClause("`fn` qualified 0.5", "test.md", 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if with != nil {
			t.Errorf("expected nil with, got %v", with)
		}
		if qual == nil || *qual != 0.5 {
			t.Errorf("expected qual=0.5, got %v", qual)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})

	t.Run("WithAndQualified", func(t *testing.T) {
		ref, with, qual, rest, err := parseUsingClause("`fn` with `term1` qualified 0.5", "test.md", 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if len(with) != 1 || with[0] != "term1" {
			t.Errorf("expected with=[term1], got %v", with)
		}
		if qual == nil || *qual != 0.5 {
			t.Errorf("expected qual=0.5, got %v", qual)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseQualifiedClause — tests parseQualifiedClause for no-digits, invalid-float, and success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseQualifiedClause(t *testing.T) {
	t.Run("NoDigits", func(t *testing.T) {
		_, _, err := parseQualifiedClause("  abc", "test.md", 1)
		if err == nil {
			t.Fatal("expected an error for missing a leading float")
		}
		if !strings.Contains(err.Error(), "expected float after 'qualified'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("InvalidFloat", func(t *testing.T) {
		_, _, err := parseQualifiedClause("--1 rest", "test.md", 2)
		if err == nil {
			t.Fatal("expected an error for an invalid float")
		}
		if !strings.Contains(err.Error(), "invalid qualified value") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		q, rest, err := parseQualifiedClause("0.75 certain", "test.md", 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q != 0.75 {
			t.Errorf("expected q=0.75, got %v", q)
		}
		if rest != " certain" {
			t.Errorf("expected rest=' certain', got %q", rest)
		}
	})
}

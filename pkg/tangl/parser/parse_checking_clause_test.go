//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCheckingClause — tests parseCheckingClause for missing-backtick error and success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseCheckingClause(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		_, _, err := parseCheckingClause("notabacktick", "test.md", 1)
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted case name after 'checking'") {
			t.Fatalf("expected error, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		name, rest, err := parseCheckingClause("`otherCase` extra", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "otherCase" {
			t.Fatalf("expected name otherCase, got %q", name)
		}
		if rest != " extra" {
			t.Fatalf("expected rest ' extra', got %q", rest)
		}
	})
}

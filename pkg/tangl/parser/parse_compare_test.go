//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompare — tests parseCompare for missing-field error, matched-operator success, and unrecognized-operator error branches
package parser

import (
	"strings"
	"testing"
)

func TestParseCompare(t *testing.T) {
	t.Run("MissingField", func(t *testing.T) {
		_, err := parseCompare("no backtick here", 1, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted field") {
			t.Fatalf("expected missing field error, got %v", err)
		}
	})

	t.Run("MatchedOperator", func(t *testing.T) {
		cmp, err := parseCompare("`amount` equals 5", 2, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Field != "amount" || cmp.Op != "equals" || cmp.Value != "5" {
			t.Fatalf("expected {amount equals 5}, got %+v", cmp)
		}
	})

	t.Run("IsEmptyOperator", func(t *testing.T) {
		cmp, err := parseCompare("`amount` is empty", 3, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Field != "amount" || cmp.Op != "is empty" || cmp.Value != "" {
			t.Fatalf("expected {amount is empty}, got %+v", cmp)
		}
	})

	t.Run("UnrecognizedOperator", func(t *testing.T) {
		_, err := parseCompare("`amount` bogus 5", 4, "test.md")
		if err == nil || !strings.Contains(err.Error(), "unrecognized comparison operator") {
			t.Fatalf("expected unrecognized operator error, got %v", err)
		}
	})
}

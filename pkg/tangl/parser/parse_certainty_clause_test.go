//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCertaintyClause — tests parseCertaintyClause for missing-if, unrecognized-operator, missing-digits, atoi-overflow, missing-percent-sign, missing-certain, trailing-text, and success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseCertaintyClause(t *testing.T) {
	t.Run("MissingIf", func(t *testing.T) {
		_, err := parseCertaintyClause("certain", "test.md", 1)
		if err == nil || !strings.Contains(err.Error(), "expected 'if' certainty clause") {
			t.Fatalf("expected missing-if error, got %v", err)
		}
	})

	t.Run("UnrecognizedOperator", func(t *testing.T) {
		_, err := parseCertaintyClause("if bogus 50% certain", "test.md", 2)
		if err == nil || !strings.Contains(err.Error(), "expected certainty operator") {
			t.Fatalf("expected operator error, got %v", err)
		}
	})

	t.Run("MissingDigits", func(t *testing.T) {
		_, err := parseCertaintyClause("if at least abc% certain", "test.md", 3)
		if err == nil || !strings.Contains(err.Error(), "expected integer percent") {
			t.Fatalf("expected integer percent error, got %v", err)
		}
	})

	t.Run("AtoiOverflow", func(t *testing.T) {
		_, err := parseCertaintyClause("if at least 99999999999999999999% certain", "test.md", 4)
		if err == nil || !strings.Contains(err.Error(), "invalid percent") {
			t.Fatalf("expected invalid percent error, got %v", err)
		}
	})

	t.Run("MissingPercentSign", func(t *testing.T) {
		_, err := parseCertaintyClause("if at least 50x certain", "test.md", 5)
		if err == nil || !strings.Contains(err.Error(), "expected '%' after percent") {
			t.Fatalf("expected percent sign error, got %v", err)
		}
	})

	t.Run("MissingCertain", func(t *testing.T) {
		_, err := parseCertaintyClause("if at least 50% only", "test.md", 6)
		if err == nil || !strings.Contains(err.Error(), "expected 'certain'") {
			t.Fatalf("expected certain error, got %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, err := parseCertaintyClause("if at least 50% certain extra", "test.md", 7)
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Fatalf("expected trailing text error, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		cert, err := parseCertaintyClause("if at least 50% certain", "test.md", 8)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cert == nil || cert.Op != "at least" || cert.Percent != 50 {
			t.Fatalf("expected {at least 50}, got %+v", cert)
		}
	})
}

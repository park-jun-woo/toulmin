//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseConstValue — tests parseConstValue for no-ref, bad-ref, trailing-text, and success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseConstValue(t *testing.T) {
	t.Run("NoAs", func(t *testing.T) {
		value, ref, err := parseConstValue("  650  ", "test.md", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref != nil {
			t.Fatalf("expected nil ref, got %v", ref)
		}
		if value != "650" {
			t.Fatalf("expected value=650, got %q", value)
		}
	})

	t.Run("BadRef", func(t *testing.T) {
		_, ref, err := parseConstValue("650 as `credit", "test.md", 2)
		if err == nil {
			t.Fatal("expected an error for an unterminated ref")
		}
		if ref != nil {
			t.Fatalf("expected nil ref, got %v", ref)
		}
		if !strings.Contains(err.Error(), "expected reference after 'as'") {
			t.Errorf("expected 'expected reference after as' error, got %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, ref, err := parseConstValue("650 as `credit`.`Threshold` extra", "test.md", 3)
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if ref != nil {
			t.Fatalf("expected nil ref, got %v", ref)
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("expected 'unexpected trailing text' error, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		value, ref, err := parseConstValue("650 as `credit`.`Threshold`", "test.md", 4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != "650" {
			t.Fatalf("expected value=650, got %q", value)
		}
		if ref == nil {
			t.Fatal("expected a non-nil ref")
		}
		if ref.Alias != "credit" || ref.Name != "Threshold" {
			t.Errorf("expected ref={Alias:credit Name:Threshold}, got %+v", ref)
		}
	})
}

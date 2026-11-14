//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRequireItem — tests parseRequireItem for field/is/required/no-type branches
package parser

import (
	"testing"
)

func TestParseRequireItem(t *testing.T) {
	t.Run("NoBacktick", func(t *testing.T) {
		it := item{Text: "field is required", Line: 1}
		req, matched, err := parseRequireItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for missing backtick-quoted field")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.Field != "" {
			t.Errorf("expected zero-value Require, got %+v", req)
		}
	})

	t.Run("NoIs", func(t *testing.T) {
		it := item{Text: "`field` needs required", Line: 2}
		req, matched, err := parseRequireItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for missing 'is'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.Field != "" {
			t.Errorf("expected zero-value Require, got %+v", req)
		}
	})

	t.Run("NoRequired", func(t *testing.T) {
		it := item{Text: "`field` is mandatory", Line: 3}
		req, matched, err := parseRequireItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for missing 'required'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.Field != "" {
			t.Errorf("expected zero-value Require, got %+v", req)
		}
	})

	t.Run("NoType", func(t *testing.T) {
		it := item{Text: "`field` is required", Line: 4}
		req, matched, err := parseRequireItem(it, "test.md")
		if !matched {
			t.Fatal("expected matched=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if req.Field != "field" || req.Type != "" {
			t.Errorf("unexpected require: %+v", req)
		}
	})
}

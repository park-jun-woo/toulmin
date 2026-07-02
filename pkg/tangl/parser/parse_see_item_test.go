//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSeeItem — tests parseSeeItem for see/alias/from/package/trailing/success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseSeeItem(t *testing.T) {
	t.Run("NoSee", func(t *testing.T) {
		it := item{Text: "look `credit` from `pkg`", Line: 1}
		_, err := parseSeeItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'see'")
		}
		if !strings.Contains(err.Error(), "expected 'see `alias` from `package`'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoAlias", func(t *testing.T) {
		it := item{Text: "see credit from `pkg`", Line: 2}
		_, err := parseSeeItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted alias")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted alias after 'see'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoFrom", func(t *testing.T) {
		it := item{Text: "see `credit` in `pkg`", Line: 3}
		_, err := parseSeeItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'from'")
		}
		if !strings.Contains(err.Error(), "expected 'from' after alias") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoPackage", func(t *testing.T) {
		it := item{Text: "see `credit` from pkg", Line: 4}
		_, err := parseSeeItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted package path")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted package path after 'from'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		it := item{Text: "see `credit` from `pkg` extra", Line: 5}
		_, err := parseSeeItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		it := item{Text: "see `credit` from `pkg`", Line: 6}
		see, err := parseSeeItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if see.Alias != "credit" || see.Package != "pkg" {
			t.Errorf("expected {Alias:credit Package:pkg}, got %+v", see)
		}
	})
}

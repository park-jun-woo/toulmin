//ff:func feature=tangl type=parser control=sequence
//ff:what TestLineIndent — tests lineIndent for space, tab, mixed, non-whitespace-terminated, and all-whitespace branches
package parser

import "testing"

func TestLineIndent(t *testing.T) {
	t.Run("Spaces", func(t *testing.T) {
		if got := lineIndent("  x"); got != 2 {
			t.Fatalf("expected 2, got %d", got)
		}
	})

	t.Run("Tab", func(t *testing.T) {
		if got := lineIndent("\tx"); got != 4 {
			t.Fatalf("expected 4, got %d", got)
		}
	})

	t.Run("MixedSpaceAndTab", func(t *testing.T) {
		if got := lineIndent(" \tx"); got != 5 {
			t.Fatalf("expected 5, got %d", got)
		}
	})

	t.Run("NoLeadingWhitespace", func(t *testing.T) {
		if got := lineIndent("x"); got != 0 {
			t.Fatalf("expected 0, got %d", got)
		}
	})

	t.Run("AllWhitespace", func(t *testing.T) {
		if got := lineIndent("   "); got != 3 {
			t.Fatalf("expected 3, got %d", got)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		if got := lineIndent(""); got != 0 {
			t.Fatalf("expected 0, got %d", got)
		}
	})
}

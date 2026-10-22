//ff:func feature=tangl type=parser control=sequence
//ff:what TestTrimStart — tests trimStart for no-leading-whitespace, leading-spaces/tabs, and all-whitespace branches
package parser

import "testing"

func TestTrimStart(t *testing.T) {
	t.Run("NoLeadingWhitespace", func(t *testing.T) {
		got := trimStart("hello")
		if got != "hello" {
			t.Fatalf("expected %q, got %q", "hello", got)
		}
	})

	t.Run("LeadingSpacesAndTabs", func(t *testing.T) {
		got := trimStart("  \t hello")
		if got != "hello" {
			t.Fatalf("expected %q, got %q", "hello", got)
		}
	})

	t.Run("AllWhitespace", func(t *testing.T) {
		got := trimStart("   \t\t")
		if got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("EmptyString", func(t *testing.T) {
		got := trimStart("")
		if got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
}

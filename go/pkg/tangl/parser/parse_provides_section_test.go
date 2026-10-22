//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseProvidesSection — tests parseProvidesSection for parseItems errors, empty sections, item errors, and successful accumulation
package parser

import "testing"

func TestParseProvidesSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. provides `a`", "3. provides `b`"}, LineOffset: 1}
		_, err := parseProvidesSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1}
		eps, err := parseProvidesSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(eps) != 0 {
			t.Fatalf("expected no endpoints, got %+v", eps)
		}
	})

	t.Run("ItemError", func(t *testing.T) {
		sec := section{Lines: []string{"- bogus"}, LineOffset: 1}
		_, err := parseProvidesSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed provides item")
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- provides `a`", "- provides `b`"}, LineOffset: 1}
		eps, err := parseProvidesSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(eps) != 2 || eps[0].Name != "a" || eps[1].Name != "b" {
			t.Fatalf("expected endpoints [a, b], got %+v", eps)
		}
	})
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSeeSection — tests parseSeeSection for parseItems errors, empty sections, item errors, and successful accumulation
package parser

import "testing"

func TestParseSeeSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. see `a` from `pkg1`", "3. see `b` from `pkg2`"}, LineOffset: 1}
		_, err := parseSeeSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1}
		sees, err := parseSeeSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sees) != 0 {
			t.Fatalf("expected no sees, got %+v", sees)
		}
	})

	t.Run("ItemError", func(t *testing.T) {
		sec := section{Lines: []string{"- bogus"}, LineOffset: 1}
		_, err := parseSeeSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed see item")
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- see `a` from `pkg1`", "- see `b` from `pkg2`"}, LineOffset: 1}
		sees, err := parseSeeSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sees) != 2 || sees[0].Alias != "a" || sees[1].Alias != "b" {
			t.Fatalf("expected sees [a, b], got %+v", sees)
		}
	})
}

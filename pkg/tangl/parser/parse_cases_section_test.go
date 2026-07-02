//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCasesSection — tests parseCasesSection for parseItems errors, empty sections, item errors, and successful accumulation
package parser

import "testing"

func TestParseCasesSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. in case of `a`", "3. in case of `b`"}, LineOffset: 1}
		_, err := parseCasesSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1}
		cases, err := parseCasesSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cases) != 0 {
			t.Fatalf("expected no cases, got %+v", cases)
		}
	})

	t.Run("ItemError", func(t *testing.T) {
		sec := section{Lines: []string{"- bogus"}, LineOffset: 1}
		_, err := parseCasesSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed case item")
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- in case of `a`", "- in case of `b`"}, LineOffset: 1}
		cases, err := parseCasesSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cases) != 2 || cases[0].Name != "a" || cases[1].Name != "b" {
			t.Fatalf("expected cases [a, b], got %+v", cases)
		}
	})
}

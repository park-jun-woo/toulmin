//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDefinitionsSection — tests parseDefinitionsSection for parseItems errors, empty sections, item errors, and successful accumulation
package parser

import "testing"

func TestParseDefinitionsSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. `a` means 1", "3. `b` means 2"}, LineOffset: 1}
		_, err := parseDefinitionsSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1}
		defs, err := parseDefinitionsSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(defs) != 0 {
			t.Fatalf("expected no definitions, got %+v", defs)
		}
	})

	t.Run("ItemError", func(t *testing.T) {
		sec := section{Lines: []string{"- bogus"}, LineOffset: 1}
		_, err := parseDefinitionsSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed definition item")
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- `a` means 1", "- `b` means 2"}, LineOffset: 1}
		defs, err := parseDefinitionsSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(defs) != 2 || defs[0].Name != "a" || defs[1].Name != "b" {
			t.Fatalf("expected definitions [a, b], got %+v", defs)
		}
	})
}

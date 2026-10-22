//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRulesSection — tests parseRulesSection for parseItems errors, empty sections, item errors, and successful accumulation
package parser

import "testing"

func TestParseRulesSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. `a` when `x` equals 1", "3. `b` when `x` equals 2"}, LineOffset: 1}
		_, err := parseRulesSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1}
		rules, err := parseRulesSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(rules) != 0 {
			t.Fatalf("expected no rules, got %+v", rules)
		}
	})

	t.Run("ItemError", func(t *testing.T) {
		sec := section{Lines: []string{"- bogus"}, LineOffset: 1}
		_, err := parseRulesSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed rule item")
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- `a` when `x` equals 1", "- `b` when `x` equals 2"}, LineOffset: 1}
		rules, err := parseRulesSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(rules) != 2 || rules[0].Name != "a" || rules[1].Name != "b" {
			t.Fatalf("expected rules [a, b], got %+v", rules)
		}
	})
}

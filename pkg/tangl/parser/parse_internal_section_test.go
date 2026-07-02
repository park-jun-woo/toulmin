//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInternalSection — tests parseInternalSection for parseItems errors, empty sections, item errors, and successful accumulation
package parser

import "testing"

func TestParseInternalSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. on start", "3. on stop"}, LineOffset: 1}
		_, err := parseInternalSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1}
		ins, err := parseInternalSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ins) != 0 {
			t.Fatalf("expected no internals, got %+v", ins)
		}
	})

	t.Run("ItemError", func(t *testing.T) {
		sec := section{Lines: []string{"- bogus"}, LineOffset: 1}
		_, err := parseInternalSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed internal item")
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- on start", "- on stop"}, LineOffset: 1}
		ins, err := parseInternalSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ins) != 2 || ins[0].Event != "start" || ins[1].Event != "stop" {
			t.Fatalf("expected internals [start, stop], got %+v", ins)
		}
	})
}

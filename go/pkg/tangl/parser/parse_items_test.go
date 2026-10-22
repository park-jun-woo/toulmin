//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseItems — tests parseItems for the checkOrderedSiblings error and success branches
package parser

import "testing"

func TestParseItems(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		_, err := parseItems([]string{"1. first", "3. second"}, 1, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("Success", func(t *testing.T) {
		items, err := parseItems([]string{"- first", "- second"}, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 2 || items[0].Text != "first" || items[1].Text != "second" {
			t.Fatalf("expected items [first, second], got %+v", items)
		}
	})
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestCheckOrderedSibling — tests checkOrderedSibling for unordered, correctly numbered, and misnumbered items
package parser

import (
	"strings"
	"testing"
)

// TestCheckOrderedSibling covers all three branches of checkOrderedSibling:
// an unordered item (returns expect unchanged), a correctly numbered ordered
// item (returns expect+1), and a misnumbered ordered item (returns an error).
func TestCheckOrderedSibling(t *testing.T) {
	t.Run("Unordered", func(t *testing.T) {
		next, err := checkOrderedSibling(item{Text: "a", Ordered: false, Line: 1}, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if next != 1 {
			t.Errorf("next = %d, want 1", next)
		}
	})

	t.Run("CorrectlyNumbered", func(t *testing.T) {
		next, err := checkOrderedSibling(item{Text: "a", Ordered: true, Number: 1, Line: 1}, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if next != 2 {
			t.Errorf("next = %d, want 2", next)
		}
	})

	t.Run("Misnumbered", func(t *testing.T) {
		_, err := checkOrderedSibling(item{Text: "a", Ordered: true, Number: 3, Line: 2}, 1, "test.md")
		if err == nil {
			t.Fatal("expected an error for a misnumbered ordered item")
		}
		if !strings.Contains(err.Error(), "ordered list item numbered 3, expected 1") {
			t.Errorf("expected misnumbering error, got %v", err)
		}
	})
}

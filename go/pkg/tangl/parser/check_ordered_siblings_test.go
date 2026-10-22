//ff:func feature=tangl type=parser control=sequence
//ff:what TestCheckOrderedSiblings — tests checkOrderedSiblings for empty input, unordered items, correctly numbered items, misnumbered items, and child-error propagation
package parser

import (
	"strings"
	"testing"
)

func TestCheckOrderedSiblings(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		if err := checkOrderedSiblings(nil, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("UnorderedItems", func(t *testing.T) {
		items := []item{
			{Text: "a", Ordered: false, Line: 1},
			{Text: "b", Ordered: false, Line: 2},
		}
		if err := checkOrderedSiblings(items, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("CorrectlyNumbered", func(t *testing.T) {
		items := []item{
			{Text: "a", Ordered: true, Number: 1, Line: 1},
			{Text: "b", Ordered: true, Number: 2, Line: 2},
		}
		if err := checkOrderedSiblings(items, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Misnumbered", func(t *testing.T) {
		items := []item{
			{Text: "a", Ordered: true, Number: 1, Line: 1},
			{Text: "b", Ordered: true, Number: 3, Line: 2},
		}
		err := checkOrderedSiblings(items, "test.md")
		if err == nil {
			t.Fatal("expected an error for a misnumbered ordered item")
		}
		if !strings.Contains(err.Error(), "ordered list item numbered 3, expected 2") {
			t.Errorf("expected misnumbering error, got %v", err)
		}
	})

	t.Run("ChildErrorPropagates", func(t *testing.T) {
		items := []item{
			{
				Text: "parent", Ordered: false, Line: 1,
				Children: []item{
					{Text: "child", Ordered: true, Number: 2, Line: 2},
				},
			},
		}
		err := checkOrderedSiblings(items, "test.md")
		if err == nil {
			t.Fatal("expected an error propagated from a child item")
		}
		if !strings.Contains(err.Error(), "ordered list item numbered 2, expected 1") {
			t.Errorf("expected child misnumbering error, got %v", err)
		}
	})
}

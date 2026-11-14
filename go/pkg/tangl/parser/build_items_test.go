//ff:func feature=tangl type=parser control=sequence
//ff:what TestBuildItems — tests buildItems for empty input, non-list lines, nested children, ordered items, struck items, and mismatched-indent skip
package parser

import "testing"

func TestBuildItems(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		if items := buildItems(nil); items != nil {
			t.Fatalf("expected nil for empty entries, got %+v", items)
		}
	})

	t.Run("NonListLineIgnored", func(t *testing.T) {
		entries := []lineEntry{
			{Indent: 0, Text: "not a list item", Line: 1},
		}
		items := buildItems(entries)
		if len(items) != 0 {
			t.Fatalf("expected no items for a non-list line, got %+v", items)
		}
	})

	t.Run("NestedChildren", func(t *testing.T) {
		entries := []lineEntry{
			{Indent: 0, Text: "- parent", Line: 1},
			{Indent: 2, Text: "- child", Line: 2},
			{Indent: 0, Text: "- sibling", Line: 3},
		}
		items := buildItems(entries)
		if len(items) != 2 {
			t.Fatalf("expected 2 top-level items, got %+v", items)
		}
		if items[0].Text != "parent" {
			t.Fatalf("expected first item text 'parent', got %q", items[0].Text)
		}
		if len(items[0].Children) != 1 || items[0].Children[0].Text != "child" {
			t.Fatalf("expected parent to have one child 'child', got %+v", items[0].Children)
		}
		if items[1].Text != "sibling" {
			t.Fatalf("expected second item text 'sibling', got %q", items[1].Text)
		}
	})

	t.Run("OrderedItem", func(t *testing.T) {
		entries := []lineEntry{
			{Indent: 0, Text: "1. first", Line: 1},
		}
		items := buildItems(entries)
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %+v", items)
		}
		if !items[0].Ordered || items[0].Number != 1 || items[0].Text != "first" {
			t.Fatalf("expected ordered item {1, first}, got %+v", items[0])
		}
	})

	t.Run("StruckItem", func(t *testing.T) {
		entries := []lineEntry{
			{Indent: 0, Text: "- ~~ignored~~", Line: 1},
		}
		items := buildItems(entries)
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %+v", items)
		}
		if !items[0].Struck || items[0].Text != "ignored" {
			t.Fatalf("expected struck item with text 'ignored', got %+v", items[0])
		}
	})

	t.Run("MismatchedIndentSkipped", func(t *testing.T) {
		entries := []lineEntry{
			{Indent: 2, Text: "- a", Line: 1},
			{Indent: 0, Text: "- b", Line: 2},
		}
		items := buildItems(entries)
		if len(items) != 1 {
			t.Fatalf("expected only 1 item (mismatched-indent entry skipped), got %+v", items)
		}
		if items[0].Text != "a" {
			t.Fatalf("expected surviving item text 'a', got %q", items[0].Text)
		}
	})
}

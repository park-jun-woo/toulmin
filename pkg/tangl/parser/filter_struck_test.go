//ff:func feature=tangl type=parser control=sequence
//ff:what TestFilterStruck — tests filterStruck for struck-item removal, kept-item retention, and nested children filtering
package parser

import "testing"

func TestFilterStruck(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		if items := filterStruck(nil); len(items) != 0 {
			t.Fatalf("expected no items, got %+v", items)
		}
	})

	t.Run("RemovesStruckKeepsOthers", func(t *testing.T) {
		items := []item{
			{Text: "kept", Struck: false},
			{Text: "removed", Struck: true},
		}
		got := filterStruck(items)
		if len(got) != 1 || got[0].Text != "kept" {
			t.Fatalf("expected only 'kept' item, got %+v", got)
		}
	})

	t.Run("FiltersNestedChildren", func(t *testing.T) {
		items := []item{
			{
				Text: "parent", Struck: false,
				Children: []item{
					{Text: "childKept", Struck: false},
					{Text: "childRemoved", Struck: true},
				},
			},
		}
		got := filterStruck(items)
		if len(got) != 1 {
			t.Fatalf("expected 1 top-level item, got %+v", got)
		}
		if len(got[0].Children) != 1 || got[0].Children[0].Text != "childKept" {
			t.Fatalf("expected only 'childKept' child, got %+v", got[0].Children)
		}
	})
}

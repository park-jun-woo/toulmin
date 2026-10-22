//ff:func feature=tangl type=parser control=sequence
//ff:what TestBuildItemTree — tests buildItemTree for empty input, flat siblings, nested children, and orphaned lines
package scan

import "testing"

func TestBuildItemTree(t *testing.T) {
	// Empty input yields a nil tree.
	if got := buildItemTree(nil); got != nil {
		t.Errorf("buildItemTree(nil) = %+v, want nil", got)
	}

	// Two flat siblings: neither has children (inner loop never fires).
	siblings := []Line{
		{Span: Span{Off: 0, Len: 6}, Text: "- a", No: 1},
		{Span: Span{Off: 6, Len: 6}, Text: "- b", No: 2},
	}
	items := buildItemTree(siblings)
	if len(items) != 2 {
		t.Fatalf("siblings: got %d items, want 2: %+v", len(items), items)
	}
	if items[0].Text != "a" || items[0].Indent != 0 || len(items[0].Children) != 0 {
		t.Errorf("siblings[0] = %+v", items[0])
	}
	if items[1].Text != "b" || items[1].Indent != 0 || len(items[1].Children) != 0 {
		t.Errorf("siblings[1] = %+v", items[1])
	}

	// A grandchild forces the inner loop to fire at every nesting level.
	nested := []Line{
		{Span: Span{Off: 0, Len: 6}, Text: "- top", No: 1},
		{Span: Span{Off: 6, Len: 10}, Text: "  - child", No: 2},
		{Span: Span{Off: 16, Len: 16}, Text: "    - grandchild", No: 3},
	}
	top := buildItemTree(nested)
	if len(top) != 1 || top[0].Text != "top" {
		t.Fatalf("nested top = %+v, want one item Text=top", top)
	}
	if len(top[0].Children) != 1 || top[0].Children[0].Text != "child" {
		t.Fatalf("nested top[0].Children = %+v, want one item Text=child", top[0].Children)
	}
	grandchild := top[0].Children[0].Children
	if len(grandchild) != 1 || grandchild[0].Text != "grandchild" {
		t.Errorf("nested grandchild = %+v, want one item Text=grandchild", grandchild)
	}

	// A line indented less than the local base (but still part of the slice)
	// is neither an item at this level nor a child: it is skipped outright.
	orphaned := []Line{
		{Span: Span{Off: 0, Len: 18}, Text: "    - nested deep", No: 1},
		{Span: Span{Off: 18, Len: 14}, Text: "  - shallower", No: 2},
	}
	skip := buildItemTree(orphaned)
	if len(skip) != 1 || skip[0].Text != "nested deep" {
		t.Fatalf("orphaned = %+v, want one item Text=\"nested deep\"", skip)
	}
	if len(skip[0].Children) != 0 {
		t.Errorf("orphaned[0].Children = %+v, want none (shallower line is dropped)", skip[0].Children)
	}
}

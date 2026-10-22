//ff:func feature=tangl type=parser control=sequence
//ff:what TestScanItems — tests ScanItems for nesting, ordered index, strikethrough, and interstitial prose
package scan

import "testing"

func TestScanItems(t *testing.T) {
	src := []byte("- top\n" +
		"  - child a\n" +
		"  - ~~struck~~\n" +
		"\n" +
		"prose line\n" +
		"1. ordered one\n" +
		"2. ordered two\n")
	lines := ScanLines(src)
	items, prose := ScanItems(lines)

	// Three top-level items: the "- top" tree and two ordered siblings.
	if len(items) != 3 {
		t.Fatalf("got %d top-level items, want 3: %+v", len(items), items)
	}

	top := items[0]
	if top.Text != "top" || top.Ordered || top.Indent != 0 {
		t.Errorf("top = %+v", top)
	}
	if len(top.Children) != 2 {
		t.Fatalf("top.Children = %d, want 2", len(top.Children))
	}
	if top.Children[0].Text != "child a" || top.Children[0].Indent != 2 {
		t.Errorf("child a = %+v", top.Children[0])
	}
	// Strikethrough item keeps its span but is flagged and unwrapped.
	struck := top.Children[1]
	if !struck.Struck || struck.Text != "struck" {
		t.Errorf("struck child = %+v, want Struck && Text=struck", struck)
	}

	// Ordered markers preserve their number.
	if !items[1].Ordered || items[1].Index != 1 || items[1].Text != "ordered one" {
		t.Errorf("items[1] = %+v", items[1])
	}
	if !items[2].Ordered || items[2].Index != 2 {
		t.Errorf("items[2] = %+v", items[2])
	}

	// The blank line and the non-list line become prose.
	if len(prose) != 2 {
		t.Fatalf("prose = %d lines, want 2: %+v", len(prose), prose)
	}
	if prose[0].Text != "" || prose[1].Text != "prose line" {
		t.Errorf("prose = %q / %q", prose[0].Text, prose[1].Text)
	}

	// Item.Span covers the whole line including the marker and newline.
	if string(src[top.Span.Off:top.Span.Off+top.Span.Len]) != "- top\n" {
		t.Errorf("top.Span bytes = %q", string(src[top.Span.Off:top.Span.Off+top.Span.Len]))
	}
}

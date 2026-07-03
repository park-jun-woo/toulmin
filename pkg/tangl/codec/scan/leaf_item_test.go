//ff:func feature=tangl type=parser control=sequence
//ff:what TestLeafItem — tests leafItem for unordered, ordered, and strikethrough markers, and span/no/indent passthrough
package scan

import "testing"

func TestLeafItem(t *testing.T) {
	// Unordered marker: Ordered/Index stay zero-valued.
	unordered := leafItem(Line{Span: Span{Off: 4, Len: 8}, Text: "- plain", No: 3}, 2)
	if unordered.Text != "plain" || unordered.Ordered || unordered.Index != 0 {
		t.Errorf("unordered = %+v", unordered)
	}
	if unordered.Span.Off != 4 || unordered.Span.Len != 8 || unordered.No != 3 || unordered.Indent != 2 {
		t.Errorf("unordered span/no/indent = %+v", unordered)
	}
	if unordered.Struck {
		t.Errorf("unordered.Struck = true, want false")
	}

	// Ordered marker: Ordered is set and Index carries the parsed number.
	ordered := leafItem(Line{Span: Span{Off: 0, Len: 10}, Text: "3. third", No: 5}, 0)
	if !ordered.Ordered || ordered.Index != 3 || ordered.Text != "third" {
		t.Errorf("ordered = %+v", ordered)
	}

	// Strikethrough content is unwrapped and flagged.
	struck := leafItem(Line{Span: Span{Off: 0, Len: 12}, Text: "- ~~gone~~", No: 7}, 4)
	if !struck.Struck || struck.Text != "gone" {
		t.Errorf("struck = %+v, want Struck && Text=gone", struck)
	}

	// Leading/trailing whitespace on the raw line is trimmed before parsing.
	padded := leafItem(Line{Span: Span{Off: 0, Len: 6}, Text: "  - x  ", No: 9}, 2)
	if padded.Text != "x" || padded.Ordered {
		t.Errorf("padded = %+v", padded)
	}
}

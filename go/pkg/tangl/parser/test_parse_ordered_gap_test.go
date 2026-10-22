//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseOrderedGap — a numbering gap in an ordered list is a parser error
package parser

import "testing"

// TestParseOrderedGap verifies that an ordered ("N. ") list whose numbers
// skip (1, 3, ...) is rejected.
func TestParseOrderedGap(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n## tangl:Cases\n- in case of `x`\n  1. `a` is required\n  3. `b` is required\n"
	_, err := ParseSource(src, "ordered_gap.md")
	if err == nil {
		t.Fatal("expected error for ordered list numbering gap, got nil")
	}
}

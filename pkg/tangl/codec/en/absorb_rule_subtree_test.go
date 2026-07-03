//ff:func feature=tangl type=parser control=sequence
//ff:what TestAbsorbRuleSubtree — a Rules head absorbs its nested condition subtree through the last line
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// TestAbsorbRuleSubtree checks that a rule head's absorbed span reaches the end of
// its deepest condition line: the ordered head with an either/nested subtree ends
// at line 3, and the second head with a single nested line ends at line 5. The
// span end is taken from the scanner line table so blank lines within a subtree
// would be covered too.
func TestAbsorbRuleSubtree(t *testing.T) {
	src := []byte("1. `r` when\n" +
		"    - either\n" +
		"      - `x` is at most 0\n" +
		"- `s` when\n" +
		"    - not `y` equals \"z\"\n")
	lines := scan.ScanLines(src)
	end := func(no int) uint32 { return lines[no-1].Span.Off + lines[no-1].Span.Len }

	items, _ := scan.ScanItems(lines)
	if len(items) != 2 {
		t.Fatalf("ScanItems top-level = %d, want 2", len(items))
	}

	endOff, lastNo := absorbRuleSubtree(items[0])
	if lastNo != 3 || endOff != end(3) {
		t.Errorf("head 0: (endOff=%d, lastNo=%d), want (%d, 3)", endOff, lastNo, end(3))
	}
	endOff, lastNo = absorbRuleSubtree(items[1])
	if lastNo != 5 || endOff != end(5) {
		t.Errorf("head 1: (endOff=%d, lastNo=%d), want (%d, 5)", endOff, lastNo, end(5))
	}
}

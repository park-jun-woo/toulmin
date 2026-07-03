//ff:func feature=tangl type=parser control=sequence
//ff:what TestCollectRuleAnchors — anchor each Rules head with its absorbed subtree span, skip struck heads
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// TestCollectRuleAnchors feeds two Rules heads, the first struck (skipped
// entirely, so its missing refByLine entry never triggers an error) and the
// second live with a nested condition line. It checks only the live head is
// anchored, with a span absorbing the nested line, then checks a live head with
// no AST ref reports the drift error.
func TestCollectRuleAnchors(t *testing.T) {
	src := []byte("1. ~~`dead` when `p`.`x`~~\n" +
		"- `r` when\n" +
		"    - `x` is at most 0\n")
	lines := scan.ScanLines(src)
	items, _ := scan.ScanItems(lines)
	if len(items) != 2 {
		t.Fatalf("ScanItems top-level = %d, want 2", len(items))
	}
	if !items[0].Struck {
		t.Fatal("expected the first Rules head to be struck")
	}

	refByLine := map[int]doc.StmtRef{2: {Sec: codebook.SectionRules, Path: []int{0}}}
	anchors, err := collectRuleAnchors(items, refByLine)
	if err != nil {
		t.Fatalf("collectRuleAnchors: %v", err)
	}
	if len(anchors) != 1 {
		t.Fatalf("anchors = %d, want 1 (struck head excluded)", len(anchors))
	}
	head := items[1]
	if anchors[0].Span.Off != head.Span.Off || anchors[0].Span.Len <= head.Span.Len {
		t.Errorf("anchor span = %+v, want an absorbed span starting at %d longer than %d", anchors[0].Span, head.Span.Off, head.Span.Len)
	}

	if _, err := collectRuleAnchors(items, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("collectRuleAnchors accepted a live Rules head with no AST ref")
	}
}

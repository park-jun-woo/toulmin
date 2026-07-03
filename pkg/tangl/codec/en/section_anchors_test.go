//ff:func feature=tangl type=parser control=sequence
//ff:what TestSectionAnchors — Rules dispatches to collectRuleAnchors, everything else to collectItemAnchors
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// TestSectionAnchors distinguishes the two dispatch targets by their differing
// span behavior: a Rules head's Block span absorbs its whole nested condition
// subtree (collectRuleAnchors), while a Cases item's Block span covers only its
// own line (collectItemAnchors), even though both source snippets have the same
// shape (a head with one nested child).
func TestSectionAnchors(t *testing.T) {
	src := []byte("1. `r` when\n" +
		"    - `x` is at most 0\n")
	lines := scan.ScanLines(src)
	items, _ := scan.ScanItems(lines)
	if len(items) != 1 {
		t.Fatalf("ScanItems = %d items, want 1", len(items))
	}
	headSpan := items[0].Span.Doc()
	refByLine := map[int]doc.StmtRef{1: {Sec: codebook.SectionRules, Path: []int{0}}}

	rulesAnchors, err := sectionAnchors(items, codebook.SectionRules, refByLine)
	if err != nil {
		t.Fatalf("sectionAnchors(Rules): %v", err)
	}
	if len(rulesAnchors) != 1 || rulesAnchors[0].Span.Len <= headSpan.Len {
		t.Errorf("Rules dispatch did not absorb the nested subtree: %+v (head span %+v)", rulesAnchors, headSpan)
	}

	refByLine = map[int]doc.StmtRef{
		1: {Sec: codebook.SectionCases, Path: []int{0}},
		2: {Sec: codebook.SectionCases, Path: []int{0, 1, 0}},
	}
	itemAnchors, err := sectionAnchors(items, codebook.SectionCases, refByLine)
	if err != nil {
		t.Fatalf("sectionAnchors(Cases): %v", err)
	}
	if len(itemAnchors) != 2 {
		t.Fatalf("Cases dispatch = %d anchors, want 2 (one per item, not absorbed)", len(itemAnchors))
	}
	if itemAnchors[0].Span.Len != headSpan.Len {
		t.Errorf("Cases dispatch head span = %+v, want the single-line span %+v", itemAnchors[0].Span, headSpan)
	}
}

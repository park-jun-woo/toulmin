//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestCollectItemAnchors — flatten a struck-skipping item tree into single-line StmtBlocks
package en

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// TestCollectItemAnchors feeds a Cases section with a live head, a live nested
// child, and a struck nested sibling (skipped along with its own child). It
// checks the struck item and its descendant never reach itemRef (both are absent
// from refByLine, which would otherwise error), and that a live item lacking a
// ref reports the drift error.
func TestCollectItemAnchors(t *testing.T) {
	src := []byte("- in case of `c`\n" +
		"  - `n` is a node\n" +
		"  - ~~`dead` is a node~~\n" +
		"    - `also dead`\n")
	lines := scan.ScanLines(src)
	items, _ := scan.ScanItems(lines)
	if len(items) != 1 {
		t.Fatalf("ScanItems top-level = %d, want 1", len(items))
	}

	refByLine := map[int]doc.StmtRef{
		1: {Sec: toulmin.SectionCases, Path: []int{0}},
		2: {Sec: toulmin.SectionCases, Path: []int{0, 1, 0}},
	}
	anchors, err := collectItemAnchors(items, toulmin.SectionCases, refByLine)
	if err != nil {
		t.Fatalf("collectItemAnchors: %v", err)
	}
	if len(anchors) != 2 {
		t.Fatalf("anchors = %d, want 2 (struck item and its child excluded)", len(anchors))
	}
	for _, a := range anchors {
		if a.Kind != doc.StmtBlock {
			t.Errorf("anchor kind = %v, want StmtBlock", a.Kind)
		}
	}

	// A live item with no entry in refByLine is a scanner/parser drift.
	if _, err := collectItemAnchors(items, toulmin.SectionCases, map[int]doc.StmtRef{}); err == nil {
		t.Fatal("collectItemAnchors accepted a live item with no AST ref")
	}
}

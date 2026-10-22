//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestMergeProse — unsorted anchors are sorted and gaps become tiling ProseBlocks
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestMergeProse feeds anchors out of order with a leading gap, an interior gap,
// and trailing bytes, and asserts mergeProse sorts them and fills every gap so the
// result tiles [0, total). The interior anchors are preserved between the prose.
func TestMergeProse(t *testing.T) {
	total := 15
	anchors := []doc.Block{
		{Kind: doc.StmtBlock, Span: doc.Span{Off: 10, Len: 2}},
		{Kind: doc.SectionBlock, Span: doc.Span{Off: 3, Len: 3}},
	}
	got := mergeProse(anchors, total)

	want := []doc.Block{
		{Kind: doc.ProseBlock, Span: doc.Span{Off: 0, Len: 3}},
		{Kind: doc.SectionBlock, Span: doc.Span{Off: 3, Len: 3}},
		{Kind: doc.ProseBlock, Span: doc.Span{Off: 6, Len: 4}},
		{Kind: doc.StmtBlock, Span: doc.Span{Off: 10, Len: 2}},
		{Kind: doc.ProseBlock, Span: doc.Span{Off: 12, Len: 3}},
	}
	if len(got) != len(want) {
		t.Fatalf("mergeProse len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i].Kind != want[i].Kind || got[i].Span != want[i].Span {
			t.Errorf("block %d = {%d %+v}, want {%d %+v}", i, got[i].Kind, got[i].Span, want[i].Kind, want[i].Span)
		}
	}
	// The assembled tiling must satisfy the doc-layer gate.
	if err := doc.VerifyTiling(&doc.Doc{Source: make([]byte, total), Blocks: got}); err != nil {
		t.Errorf("VerifyTiling: %v", err)
	}
}

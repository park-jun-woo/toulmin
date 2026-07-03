//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestMergeProseKo — unsorted anchors are sorted and gaps become tiling ProseBlocks
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestMergeProseKo feeds anchors out of order with a leading gap, an interior gap,
// and trailing bytes, and asserts mergeProseKo sorts them and fills every gap so
// the result tiles [0, total). It also feeds adjacent anchors with no gap and no
// trailing bytes, so neither gap-emitting if takes its branch.
func TestMergeProseKo(t *testing.T) {
	total := 15
	anchors := []doc.Block{
		{Kind: doc.StmtBlock, Span: doc.Span{Off: 10, Len: 2}},
		{Kind: doc.SectionBlock, Span: doc.Span{Off: 3, Len: 3}},
	}
	got := mergeProseKo(anchors, total)

	want := []doc.Block{
		{Kind: doc.ProseBlock, Span: doc.Span{Off: 0, Len: 3}},
		{Kind: doc.SectionBlock, Span: doc.Span{Off: 3, Len: 3}},
		{Kind: doc.ProseBlock, Span: doc.Span{Off: 6, Len: 4}},
		{Kind: doc.StmtBlock, Span: doc.Span{Off: 10, Len: 2}},
		{Kind: doc.ProseBlock, Span: doc.Span{Off: 12, Len: 3}},
	}
	if len(got) != len(want) {
		t.Fatalf("mergeProseKo len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i].Kind != want[i].Kind || got[i].Span != want[i].Span {
			t.Errorf("block %d = {%d %+v}, want {%d %+v}", i, got[i].Kind, got[i].Span, want[i].Kind, want[i].Span)
		}
	}
	if err := doc.VerifyTiling(&doc.Doc{Source: make([]byte, total), Blocks: got}); err != nil {
		t.Errorf("VerifyTiling: %v", err)
	}

	// No gap before the first anchor and none after the last: neither
	// gap-emitting `if` body runs on this call.
	noGap := mergeProseKo([]doc.Block{{Kind: doc.StmtBlock, Span: doc.Span{Off: 0, Len: 5}}}, 5)
	if len(noGap) != 1 || noGap[0].Span != (doc.Span{Off: 0, Len: 5}) {
		t.Fatalf("noGap = %+v", noGap)
	}
}

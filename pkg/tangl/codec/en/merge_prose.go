//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what mergeProse — fill the gaps between sorted anchor blocks with verbatim ProseBlocks
package en

import (
	"sort"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// mergeProse assembles the final stream-order tiling: it sorts the anchor blocks
// (section headers and statements) by offset and emits one ProseBlock per gap —
// leading preamble, inter-section prose, blank lines, struck items, trailing bytes
// — so the returned Blocks tile [0, total) with no gap or overlap. Each residual
// run between two adjacent anchors becomes a single merged ProseBlock.
func mergeProse(anchors []doc.Block, total int) []doc.Block {
	sort.SliceStable(anchors, func(i, j int) bool {
		return anchors[i].Span.Off < anchors[j].Span.Off
	})
	var blocks []doc.Block
	var cursor uint32
	for _, a := range anchors {
		if a.Span.Off > cursor {
			blocks = append(blocks, doc.Block{Kind: doc.ProseBlock, Span: doc.Span{Off: cursor, Len: a.Span.Off - cursor}})
		}
		blocks = append(blocks, a)
		cursor = a.Span.Off + a.Span.Len
	}
	if cursor < uint32(total) {
		blocks = append(blocks, doc.Block{Kind: doc.ProseBlock, Span: doc.Span{Off: cursor, Len: uint32(total) - cursor}})
	}
	return blocks
}

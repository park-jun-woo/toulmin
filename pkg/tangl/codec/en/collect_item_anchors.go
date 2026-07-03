//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what collectItemAnchors — anchor every non-struck list item (and its descendants) as a StmtBlock
package en

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// collectItemAnchors flattens a section's item tree (an explicit stack, order
// irrelevant since mergeProse re-sorts) and emits one single-line StmtBlock per
// non-struck item, its ref from itemRef. Struck items and their subtrees are left
// uncovered for mergeProse to fill as prose.
func collectItemAnchors(items []scan.Item, code codebook.SectionCode, refByLine map[int]doc.StmtRef) ([]doc.Block, error) {
	var anchors []doc.Block
	stack := append([]scan.Item(nil), items...)
	for len(stack) > 0 {
		it := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if it.Struck {
			continue
		}
		ref, err := itemRef(code, it.No, refByLine)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: ref})
		stack = append(stack, it.Children...)
	}
	return anchors, nil
}

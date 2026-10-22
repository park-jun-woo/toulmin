//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what absorbRuleSubtree — compute the end offset and last line of a Rules head plus its condition subtree
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"

// absorbRuleSubtree returns the past-the-end byte offset and last source line of a
// Rules inline-rule head together with its whole nested condition subtree (README
// decision 4). Head and descendants occupy contiguous source lines, so the span
// from the head's offset to endOff covers the entire subtree — including any blank
// lines interleaved between condition items — as one StmtBlock.
func absorbRuleSubtree(head scan.Item) (endOff uint32, lastNo int) {
	endOff = head.Span.Off + head.Span.Len
	lastNo = head.No
	for _, child := range head.Children {
		e, n := absorbRuleSubtree(child)
		if e > endOff {
			endOff = e
		}
		if n > lastNo {
			lastNo = n
		}
	}
	return endOff, lastNo
}

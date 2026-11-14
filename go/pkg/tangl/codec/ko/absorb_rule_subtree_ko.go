//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what absorbRuleSubtreeKo — compute the past-the-end offset of a Rules head plus its condition subtree
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"

// absorbRuleSubtreeKo returns the past-the-end byte offset of a Rules inline-rule
// head together with its whole nested condition subtree (README contract §2 — the
// subtree is one StmtBlock). Head and descendants occupy contiguous source lines,
// so the span from the head's offset to this end covers the entire subtree. It
// mirrors pkg/tangl/codec/en.absorbRuleSubtree.
func absorbRuleSubtreeKo(head scan.Item) uint32 {
	endOff := head.Span.Off + head.Span.Len
	for _, child := range head.Children {
		if e := absorbRuleSubtreeKo(child); e > endOff {
			endOff = e
		}
	}
	return endOff
}

//ff:func feature=tangl type=parser control=sequence
//ff:what sectionAnchors — dispatch a section's statement anchoring by Rules vs. everything else
package en

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// sectionAnchors returns the StmtBlocks of one section: Rules heads absorb their
// condition subtrees (collectRuleAnchors), every other section anchors one Block per
// non-struck list item (collectItemAnchors).
func sectionAnchors(items []scan.Item, code codebook.SectionCode, refByLine map[int]doc.StmtRef) ([]doc.Block, error) {
	if code == codebook.SectionRules {
		return collectRuleAnchors(items, refByLine)
	}
	return collectItemAnchors(items, code, refByLine)
}

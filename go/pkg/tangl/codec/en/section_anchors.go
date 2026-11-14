//ff:func feature=tangl type=parser control=sequence
//ff:what sectionAnchors — dispatch a section's statement anchoring by Rules vs. everything else
package en

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// sectionAnchors returns the StmtBlocks of one section: Rules heads absorb their
// condition subtrees (collectRuleAnchors), every other section anchors one Block per
// non-struck list item (collectItemAnchors).
func sectionAnchors(items []scan.Item, code toulmin.SectionCode, refByLine map[int]doc.StmtRef) ([]doc.Block, error) {
	if code == toulmin.SectionRules {
		return collectRuleAnchors(items, refByLine)
	}
	return collectItemAnchors(items, code, refByLine)
}

//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what collectAnchors — build section-header and statement anchor Blocks from the scanned sections
package en

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// collectAnchors produces one SectionBlock per `## tangl:` heading and delegates
// the section's statement anchors to collectRuleAnchors (a Rules head absorbs its
// condition subtree into one StmtBlock) or collectItemAnchors (every other section,
// one StmtBlock per non-struck list item). Struck items are left for mergeProse to
// cover as prose. An unknown section name is a scanner/parser drift error.
func collectAnchors(secs []scan.Section, refByLine map[int]doc.StmtRef) ([]doc.Block, error) {
	var anchors []doc.Block
	for _, sec := range secs {
		code, ok := sectionCode(sec.Name)
		if !ok {
			return nil, fmt.Errorf("tangl/codec/en: unknown tangl section %q", sec.Name)
		}
		anchors = append(anchors, doc.Block{
			Kind: doc.SectionBlock,
			Span: sec.HeaderSpan.Doc(),
			Ref:  doc.StmtRef{Sec: code},
		})
		items, _ := scan.ScanItems(sec.Lines)
		sub, err := sectionAnchors(items, code, refByLine)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, sub...)
	}
	return anchors, nil
}

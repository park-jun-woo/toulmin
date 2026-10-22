//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what collectAnchorsKo — walk the KO sections, parsing the AST and building header/statement anchors
package ko

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// collectAnchorsKo walks the scanned `## tangl:` sections in document order, mapping
// each Korean heading to its canonical name (koSectionNames) and codebook code. It
// emits a SectionBlock per heading, enforces the required section order (codes rise
// 1..7), and delegates item parsing/anchoring to sectionAnchorsKo, which fills the
// Document. An unknown heading or an out-of-order section is a drift error.
func collectAnchorsKo(a *ast.Document, secs []scan.Section, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	maxCode := toulmin.SectionCode(0)
	for _, sec := range secs {
		canon, ok := koSectionNames(sec.Name)
		if !ok {
			return nil, fmt.Errorf("%s: unknown tangl section %q", path, sec.Name)
		}
		code, _ := sectionCodeKo(canon)
		if code < maxCode {
			return nil, fmt.Errorf("%s: tangl:%s appears out of order", path, sec.Name)
		}
		maxCode = code
		anchors = append(anchors, doc.Block{Kind: doc.SectionBlock, Span: sec.HeaderSpan.Doc(), Ref: doc.StmtRef{Sec: code}})
		items, _ := scan.ScanItems(sec.Lines)
		sub, err := sectionAnchorsKo(a, code, items, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, sub...)
	}
	return anchors, nil
}

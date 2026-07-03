//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what collectRuleAnchors — anchor each Rules head as one StmtBlock covering its condition subtree
package en

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// collectRuleAnchors emits one StmtBlock per non-struck Rules head, its span
// absorbing the whole nested condition subtree (README decision 4). A head with no
// AST ref is a scanner/parser drift reported as error.
func collectRuleAnchors(items []scan.Item, refByLine map[int]doc.StmtRef) ([]doc.Block, error) {
	var anchors []doc.Block
	for _, head := range items {
		if head.Struck {
			continue
		}
		ref, ok := refByLine[head.No]
		if !ok {
			return nil, fmt.Errorf("tangl/codec/en: no AST ref for Rules item at line %d", head.No)
		}
		endOff, _ := absorbRuleSubtree(head)
		anchors = append(anchors, doc.Block{
			Kind: doc.StmtBlock,
			Span: doc.Span{Off: head.Span.Off, Len: endOff - head.Span.Off},
			Ref:  ref,
		})
	}
	return anchors, nil
}

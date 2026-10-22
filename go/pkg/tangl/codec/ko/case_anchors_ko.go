//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what caseAnchorsKo — parse the 사안 section, anchoring each case head and delegating its body
package ko

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// caseAnchorsKo parses each non-struck 사안 item into an ast.Case: the head anchors
// as {Cases, [ci]} and caseBodyAnchorsKo walks the body children. Struck items are
// left for mergeProseKo to cover as prose.
func caseAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	ci := 0
	for _, it := range items {
		if it.Struck {
			continue
		}
		c, err := parseCaseKo(it, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{ci}}})
		body, err := caseBodyAnchorsKo(&c, ci, it.Children, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, body...)
		a.Cases = append(a.Cases, c)
		ci++
	}
	return anchors, nil
}

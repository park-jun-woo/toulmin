//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what caseBodyAnchorsKo — dispatch case-body children and anchor each as a {Cases,[ci,sk,si]} slot
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// caseBodyAnchorsKo walks a case's non-struck body children, dispatching each via
// applyCaseChildKo (which appends onto c) and anchoring it as {Cases, [ci, slotKind,
// slotIdx]} with per-kind counters (0 Require · 1 Node · 2 Attack · 3 Exec · 4 Supersede).
func caseBodyAnchorsKo(c *ast.Case, ci int, children []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	counts := [5]int{}
	for _, ch := range children {
		if ch.Struck {
			continue
		}
		sk, err := applyCaseChildKo(c, ch, path)
		if err != nil {
			return nil, err
		}
		si := counts[sk]
		counts[sk]++
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: ch.Span.Doc(), Ref: doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, sk, si}}})
	}
	return anchors, nil
}

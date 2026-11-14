//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what internalAnchorsKo — parse the 내부 section, anchoring each trigger head and delegating its body
package ko

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// internalAnchorsKo parses each non-struck 내부 item into an ast.Internal: the
// trigger head anchors as {Internal, [ii]} and internalBodyAnchorsKo walks the
// run/check body children. Struck items become prose via mergeProseKo.
func internalAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	ii := 0
	for _, it := range items {
		if it.Struck {
			continue
		}
		in, err := parseInternalHeadKo(it, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{ii}}})
		body, err := internalBodyAnchorsKo(&in, ii, it.Children, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, body...)
		a.Internals = append(a.Internals, in)
		ii++
	}
	return anchors, nil
}

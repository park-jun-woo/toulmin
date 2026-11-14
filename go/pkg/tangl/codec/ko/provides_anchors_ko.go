//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what providesAnchorsKo — parse the 업무 section, anchoring each endpoint head and delegating its body
package ko

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// providesAnchorsKo parses each non-struck 업무 item into an ast.Endpoint: the head
// anchors as {Provides, [ei]} and endpointBodyAnchorsKo walks the body children.
// Struck items become prose via mergeProseKo.
func providesAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	ei := 0
	for _, it := range items {
		if it.Struck {
			continue
		}
		ep, err := parseEndpointKo(it, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{ei}}})
		body, err := endpointBodyAnchorsKo(&ep, ei, it.Children, path)
		if err != nil {
			return nil, err
		}
		anchors = append(anchors, body...)
		a.Provides = append(a.Provides, ep)
		ei++
	}
	return anchors, nil
}

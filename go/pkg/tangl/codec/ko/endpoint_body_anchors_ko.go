//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what endpointBodyAnchorsKo — dispatch endpoint-body children (require/run/check) to {Provides,[ei,sk,si]}
package ko

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// endpointBodyAnchorsKo walks an endpoint's non-struck body children — a require
// (필수다, slotKind 0), run (진행한다, 1), or check (판정한다, 2) — appending each onto ep
// and anchoring it as {Provides, [ei, slotKind, slotIdx]} with per-kind counters.
func endpointBodyAnchorsKo(ep *ast.Endpoint, ei int, children []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	var reqN, runN, checkN int
	for _, ch := range children {
		if ch.Struck {
			continue
		}
		var sk, si int
		if req, ok, err := parseRequireKo(ch, path); err != nil {
			return nil, err
		} else if ok {
			ep.Requires = append(ep.Requires, req)
			sk, si, reqN = 0, reqN, reqN+1
		} else if name, kind, rcOK, err := parseRunCheckKo(ch, path); err != nil {
			return nil, err
		} else if rcOK && kind == 1 {
			ep.Runs = append(ep.Runs, name)
			sk, si, runN = 1, runN, runN+1
		} else if rcOK {
			ep.Checks = append(ep.Checks, name)
			sk, si, checkN = 2, checkN, checkN+1
		} else {
			return nil, errAtKo(path, ch.No, "unrecognized 업무 statement: %q", ch.Text)
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: ch.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{ei, sk, si}}})
	}
	return anchors, nil
}

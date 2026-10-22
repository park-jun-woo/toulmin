//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what internalBodyAnchorsKo — dispatch internal-body run/check children to {Internal,[ii,sk,si]}
package ko

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// internalBodyAnchorsKo walks an internal trigger's non-struck body children — a run
// (진행한다, slotKind 1) or check (판정한다, 2) — appending each onto in and anchoring it
// as {Internal, [ii, slotKind, slotIdx]} with per-kind counters (slotKind 0 unused).
func internalBodyAnchorsKo(in *ast.Internal, ii int, children []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	var runN, checkN int
	for _, ch := range children {
		if ch.Struck {
			continue
		}
		name, kind, ok, err := parseRunCheckKo(ch, path)
		if err != nil {
			return nil, err
		}
		var si int
		switch {
		case ok && kind == 1:
			si, runN = runN, runN+1
			in.Runs = append(in.Runs, name)
		case ok:
			si, checkN = checkN, checkN+1
			in.Checks = append(in.Checks, name)
		default:
			return nil, errAtKo(path, ch.No, "unrecognized 내부 statement: %q", ch.Text)
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: ch.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{ii, kind, si}}})
	}
	return anchors, nil
}

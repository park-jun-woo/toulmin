//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what seeAnchorsKo — parse the 참조 section and anchor each `see` as {See,[i]}
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// seeAnchorsKo parses each non-struck 참조 item into an ast.See appended onto
// a.Sees and anchors it as a StmtBlock with ref {See, [i]}, where i counts only the
// non-struck (parsed) declarations. Struck items become prose via mergeProseKo.
func seeAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	i := 0
	for _, it := range items {
		if it.Struck {
			continue
		}
		see, err := parseSeeKo(it, path)
		if err != nil {
			return nil, err
		}
		a.Sees = append(a.Sees, see)
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: doc.StmtRef{Sec: codebook.SectionSee, Path: []int{i}}})
		i++
	}
	return anchors, nil
}

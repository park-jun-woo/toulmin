//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what definitionAnchorsKo — parse the 정의 section, anchoring each definition head and delegating fields
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// definitionAnchorsKo parses each non-struck 정의 item into an ast.Definition and
// anchors the head line as {Definitions, [di]}; a struct definition's field children
// are anchored by definitionFieldAnchorsKo. Struck items become prose via mergeProseKo.
func definitionAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	di := 0
	for _, it := range items {
		if it.Struck {
			continue
		}
		def, err := parseDefinitionKo(it, path)
		if err != nil {
			return nil, err
		}
		a.Defs = append(a.Defs, def)
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{di}}})
		anchors = append(anchors, definitionFieldAnchorsKo(di, it.Children)...)
		di++
	}
	return anchors, nil
}

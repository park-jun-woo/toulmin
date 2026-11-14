//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what definitionFieldAnchorsKo — anchor a struct definition's non-struck field children
package ko

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// definitionFieldAnchorsKo anchors each non-struck struct-definition field child as
// {Definitions, [di, fi]} on its own line, with fi counting only the non-struck
// (parsed) fields. Struck fields become prose via mergeProseKo.
func definitionFieldAnchorsKo(di int, children []scan.Item) []doc.Block {
	var anchors []doc.Block
	fi := 0
	for _, ch := range children {
		if ch.Struck {
			continue
		}
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: ch.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{di, fi}}})
		fi++
	}
	return anchors
}

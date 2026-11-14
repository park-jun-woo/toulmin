//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what subjectAnchorsKo — parse the 주제 section and anchor its declaration as {Subject,nil}
package ko

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// subjectAnchorsKo parses each non-struck 주제 item into a.Subject and anchors it as
// a StmtBlock with the position-assigned ref {Subject, nil} (Subject has no Line).
// Struck items are left for mergeProseKo to cover as prose.
func subjectAnchorsKo(a *ast.Document, items []scan.Item, path string) ([]doc.Block, error) {
	var anchors []doc.Block
	for _, it := range items {
		if it.Struck {
			continue
		}
		name, err := parseSubjectKo(it, path)
		if err != nil {
			return nil, err
		}
		a.Subject = name
		anchors = append(anchors, doc.Block{Kind: doc.StmtBlock, Span: it.Span.Doc(), Ref: doc.StmtRef{Sec: toulmin.SectionSubject}})
	}
	return anchors, nil
}

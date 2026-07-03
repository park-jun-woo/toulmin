//ff:func feature=tangl type=util control=selection
//ff:what appendRequire — route a required-field declaration to its case or endpoint owner
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// appendRequire appends r to the Requires slice of the owner identified by the
// requires-edge src TID — a Case or an Endpoint — and returns the matching
// StmtRef (slotKind 0). The owner's recorded node type routes the field without
// any section context, restoring the document order the split AST slices lose.
func appendRequire(d *dec, owner uint16, r ast.Require) (doc.StmtRef, error) {
	switch d.nodeType[owner] {
	case codebook.NodeCase:
		ci := d.caseByTID[owner]
		fi := len(d.ast.Cases[ci].Requires)
		d.ast.Cases[ci].Requires = append(d.ast.Cases[ci].Requires, r)
		return doc.StmtRef{Sec: codebook.SectionCases, Path: []int{ci, 0, fi}}, nil
	case codebook.NodeEndpoint:
		ei := d.epByTID[owner]
		fi := len(d.ast.Provides[ei].Requires)
		d.ast.Provides[ei].Requires = append(d.ast.Provides[ei].Requires, r)
		return doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{ei, 0, fi}}, nil
	default:
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: requires edge owner %d is not a case or endpoint", owner)
	}
}

//ff:func feature=tangl type=util control=selection
//ff:what classifyCaseTID — issue TIDs for case declarations, requires and rule nodes
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// classifyCaseTID issues the TID for a Cases coordinate that declares an entity:
// the case itself (also recorded by name), a required field, or a rule node
// (also recorded under its case-scoped name). Attacks (slotKind 2) and execs
// (slotKind 3) declare nothing.
func classifyCaseTID(tt *tidTable, next *toulmin.TID, ref doc.StmtRef, elem any) {
	switch {
	case len(ref.Path) == 1:
		tt.caseByName[elem.(ast.Case).Name] = assignTID(tt, next, ref)
	case ref.Path[1] == 0:
		assignTID(tt, next, ref)
	case ref.Path[1] == 1:
		tt.nodeByCase[nodeKey(ref.Path[0], elem.(ast.Node).Name)] = assignTID(tt, next, ref)
	}
}

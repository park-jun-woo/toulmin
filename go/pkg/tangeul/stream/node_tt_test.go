//ff:func feature=tangl type=util control=sequence
//ff:what nodeTT — test helper: a TID table mapping one statement coordinate to an Edge Id
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// nodeTT returns a TID table whose only entry maps ref to id (pass id 0x10000 to
// make that node's own Edge Id overflow the 1-word width). The name maps are
// empty and non-primary ids start at 1.
func nodeTT(ref doc.StmtRef, id toulmin.TID) *tidTable {
	return &tidTable{
		byRef:       map[string]toulmin.TID{tidKey(ref): id},
		aliasByName: map[string]toulmin.TID{},
		caseByName:  map[string]toulmin.TID{},
		nodeByCase:  map[string]toulmin.TID{},
		primaryEnd:  1,
	}
}

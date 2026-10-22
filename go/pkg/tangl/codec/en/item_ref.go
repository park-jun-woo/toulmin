//ff:func feature=tangl type=parser control=sequence
//ff:what itemRef — resolve a non-Rules item's StmtRef, assigning Subject by section position
package en

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// itemRef returns the StmtRef for a non-struck item at source line no. Subject has
// no Line, so it is assigned {Subject,nil} by section position; every other section
// looks up refByLine, and a miss is a scanner/parser drift reported as error.
func itemRef(code toulmin.SectionCode, no int, refByLine map[int]doc.StmtRef) (doc.StmtRef, error) {
	if code == toulmin.SectionSubject {
		return doc.StmtRef{Sec: toulmin.SectionSubject}, nil
	}
	ref, ok := refByLine[no]
	if !ok {
		return doc.StmtRef{}, fmt.Errorf("tangl/codec/en: no AST ref for item at line %d", no)
	}
	return ref, nil
}

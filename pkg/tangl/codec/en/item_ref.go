//ff:func feature=tangl type=parser control=sequence
//ff:what itemRef — resolve a non-Rules item's StmtRef, assigning Subject by section position
package en

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// itemRef returns the StmtRef for a non-struck item at source line no. Subject has
// no Line, so it is assigned {Subject,nil} by section position; every other section
// looks up refByLine, and a miss is a scanner/parser drift reported as error.
func itemRef(code codebook.SectionCode, no int, refByLine map[int]doc.StmtRef) (doc.StmtRef, error) {
	if code == codebook.SectionSubject {
		return doc.StmtRef{Sec: codebook.SectionSubject}, nil
	}
	ref, ok := refByLine[no]
	if !ok {
		return doc.StmtRef{}, fmt.Errorf("tangl/codec/en: no AST ref for item at line %d", no)
	}
	return ref, nil
}

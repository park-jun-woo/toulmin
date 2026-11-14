//ff:func feature=tangl type=util control=sequence
//ff:what transitionCaseTT — fixture TID table with case 0 = 100 and case 1 = 200
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// transitionCaseTT returns a TID table whose case 0 is 100 and case 1 is 200;
// non-primary ids start at 300.
func transitionCaseTT() *tidTable {
	return &tidTable{byRef: map[string]toulmin.TID{
		tidKey(doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}): 100,
		tidKey(doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{1}}): 200,
	}, primaryEnd: 300}
}

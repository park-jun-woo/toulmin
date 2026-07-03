//ff:func feature=tangl type=util control=sequence
//ff:what transitionCaseTT — fixture tidTable resolving Cases paths [0]/[1] to TIDs 100/200
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// transitionCaseTT builds a tidTable resolving Cases path [0] and [1] to fixed
// TIDs, matching how transitionCase looks up a case's TID via tidKey.
func transitionCaseTT() *tidTable {
	return &tidTable{byRef: map[string]uint16{
		tidKey(doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}}): 100,
		tidKey(doc.StmtRef{Sec: codebook.SectionCases, Path: []int{1}}): 200,
	}}
}

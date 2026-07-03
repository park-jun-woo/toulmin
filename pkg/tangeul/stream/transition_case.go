//ff:func feature=tangl type=util control=sequence
//ff:what transitionCase — emit CASE_END/CASE_BEGIN controls as the walk enters or leaves a case
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// transitionCase updates the open-case state for the block about to be emitted,
// returning any CASE_END (leaving the current case) and CASE_BEGIN (entering a
// new one) controls. A block is in a case when it is a Cases StmtBlock with a
// non-empty path; its first path index is the case.
func transitionCase(tt *tidTable, st *caseState, b doc.Block) ([]word.Word, error) {
	inCase := b.Kind == doc.StmtBlock && b.Ref.Sec == codebook.SectionCases && len(b.Ref.Path) >= 1
	idx := 0
	if inCase {
		idx = b.Ref.Path[0]
	}
	var out []word.Word
	if st.open && (!inCase || idx != st.idx) {
		w, err := caseControlWords(codebook.ControlCaseEnd, st.tid)
		if err != nil {
			return nil, err
		}
		out = append(out, w...)
		st.open = false
	}
	if inCase && !st.open {
		st.tid = tt.byRef[tidKey(doc.StmtRef{Sec: codebook.SectionCases, Path: []int{idx}})]
		w, err := caseControlWords(codebook.ControlCaseBegin, st.tid)
		if err != nil {
			return nil, err
		}
		out = append(out, w...)
		st.open, st.idx = true, idx
	}
	return out, nil
}

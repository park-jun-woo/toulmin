//ff:func feature=tangl type=util control=sequence
//ff:what transitionCase — emit CASE_END/CASE_BEGIN controls as the walk enters or leaves a case
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// transitionCase updates the open-case state for the block about to be emitted,
// returning any CASE_END (leaving the current case) and CASE_BEGIN (entering a
// new one) controls. A block is in a case when it is a Cases StmtBlock with a
// non-empty path; its first path index is the case.
func transitionCase(e *enc, st *caseState, b doc.Block) []core.Word {
	inCase := b.Kind == doc.StmtBlock && b.Ref.Sec == toulmin.SectionCases && len(b.Ref.Path) >= 1
	idx := 0
	if inCase {
		idx = b.Ref.Path[0]
	}
	var out []core.Word
	if st.open && (!inCase || idx != st.idx) {
		out = append(out, controlWords(e, toulmin.ControlCaseEnd, uint64(st.tid), nil)...)
		st.open = false
	}
	if inCase && !st.open {
		st.tid = e.tt.byRef[tidKey(doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{idx}})]
		out = append(out, controlWords(e, toulmin.ControlCaseBegin, uint64(st.tid), nil)...)
		st.open, st.idx = true, idx
	}
	return out
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestTransitionCase — every open/enter combination of transitionCase's CASE_END/CASE_BEGIN emission
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestTransitionCase drives every combination of the function's two ifs: no
// case open and the block isn't a case (no-op); a case open and the block
// isn't a case, or is a different case (CASE_END only, or CASE_END+CASE_BEGIN
// together); no case open and the block is a case (CASE_BEGIN only); and a
// case open with the block the same case (no-op, case stays open). Each
// scenario asserts distinct setup/output, so the test itself selects.
func TestTransitionCase(t *testing.T) {
	caseEnd100, err := caseControlWords(codebook.ControlCaseEnd, 100)
	if err != nil {
		t.Fatalf("caseControlWords CASE_END 100: %v", err)
	}
	caseBegin100, err := caseControlWords(codebook.ControlCaseBegin, 100)
	if err != nil {
		t.Fatalf("caseControlWords CASE_BEGIN 100: %v", err)
	}
	caseBegin200, err := caseControlWords(codebook.ControlCaseBegin, 200)
	if err != nil {
		t.Fatalf("caseControlWords CASE_BEGIN 200: %v", err)
	}

	// run invokes transitionCase against the fixture tidTable, failing the
	// subtest on an unexpected error.
	run := func(t *testing.T, st *caseState, b doc.Block) []word.Word {
		t.Helper()
		out, err := transitionCase(transitionCaseTT(), st, b)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return out
	}

	t.Run("no case open, block outside a case: no-op", func(t *testing.T) {
		st := &caseState{}
		out := run(t, st, doc.Block{Kind: doc.ProseBlock})
		if len(out) != 0 {
			t.Errorf("out = %v, want empty", out)
		}
		if st.open {
			t.Errorf("st.open = true, want false")
		}
	})

	t.Run("case open, block outside a case: CASE_END only", func(t *testing.T) {
		st := &caseState{open: true, idx: 0, tid: 100}
		out := run(t, st, doc.Block{Kind: doc.ProseBlock})
		if !reflect.DeepEqual(out, caseEnd100) {
			t.Errorf("out = %v, want %v", out, caseEnd100)
		}
		if st.open {
			t.Errorf("st.open = true, want false")
		}
	})

	t.Run("case open, block is a Cases StmtBlock with empty path: CASE_END only", func(t *testing.T) {
		// len(Path) == 0 keeps inCase false even though Sec == SectionCases,
		// covering the third operand of the inCase expression independent of
		// Kind/Sec mismatches used by the other subtests.
		st := &caseState{open: true, idx: 0, tid: 100}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: codebook.SectionCases, Path: nil}})
		if !reflect.DeepEqual(out, caseEnd100) {
			t.Errorf("out = %v, want %v", out, caseEnd100)
		}
		if st.open {
			t.Errorf("st.open = true, want false")
		}
	})

	t.Run("no case open, block enters a case: CASE_BEGIN only", func(t *testing.T) {
		st := &caseState{}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}}})
		if !reflect.DeepEqual(out, caseBegin100) {
			t.Errorf("out = %v, want %v", out, caseBegin100)
		}
		if !st.open || st.idx != 0 || st.tid != 100 {
			t.Errorf("st = %+v, want open idx=0 tid=100", st)
		}
	})

	t.Run("case open, block is the same case: no-op, stays open", func(t *testing.T) {
		st := &caseState{open: true, idx: 0, tid: 100}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}}})
		if len(out) != 0 {
			t.Errorf("out = %v, want empty", out)
		}
		if !st.open || st.idx != 0 || st.tid != 100 {
			t.Errorf("st = %+v, want unchanged open idx=0 tid=100", st)
		}
	})

	t.Run("case open, block enters a different case: CASE_END then CASE_BEGIN", func(t *testing.T) {
		st := &caseState{open: true, idx: 0, tid: 100}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: codebook.SectionCases, Path: []int{1}}})
		want := append(append([]word.Word(nil), caseEnd100...), caseBegin200...)
		if !reflect.DeepEqual(out, want) {
			t.Errorf("out = %v, want %v", out, want)
		}
		if !st.open || st.idx != 1 || st.tid != 200 {
			t.Errorf("st = %+v, want open idx=1 tid=200", st)
		}
	})
}

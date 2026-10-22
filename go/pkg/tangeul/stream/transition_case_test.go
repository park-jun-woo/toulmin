//ff:func feature=tangl type=util control=sequence
//ff:what TestTransitionCase — every open/enter combination of transitionCase's CASE_END/CASE_BEGIN emission
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestTransitionCase drives every combination of the function's two ifs: no
// case open and the block isn't a case (no-op); a case open and the block
// isn't a case, or is a different case (CASE_END only, or CASE_END+CASE_BEGIN
// together); no case open and the block is a case (CASE_BEGIN only); and a
// case open with the block the same case (no-op, case stays open). Each
// scenario asserts distinct setup/output, so the test itself selects.
func TestTransitionCase(t *testing.T) {
	// run invokes transitionCase against the fixture tidTable, failing the
	// subtest on an unexpected error.
	run := func(t *testing.T, st *caseState, b doc.Block) []toulmin.ControlPacket {
		t.Helper()
		out, err := transitionCase(testEnc(transitionCaseTT()), st, b)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return controlsOf(t, out)
	}
	ctl := func(code toulmin.ControlCode, arg uint64) toulmin.ControlPacket {
		return toulmin.ControlPacket{Code: code, Arg: arg}
	}
	same := func(got []toulmin.ControlPacket, want ...toulmin.ControlPacket) bool {
		if len(got) != len(want) {
			return false
		}
		for i := range got {
			if got[i].Code != want[i].Code || got[i].Arg != want[i].Arg || got[i].ID != toulmin.TID(300+i) {
				return false
			}
		}
		return true
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
		if !same(out, ctl(toulmin.ControlCaseEnd, 100)) {
			t.Errorf("out = %+v, want CASE-END 100", out)
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
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: toulmin.SectionCases, Path: nil}})
		if !same(out, ctl(toulmin.ControlCaseEnd, 100)) {
			t.Errorf("out = %+v, want CASE-END 100", out)
		}
		if st.open {
			t.Errorf("st.open = true, want false")
		}
	})

	t.Run("no case open, block enters a case: CASE_BEGIN only", func(t *testing.T) {
		st := &caseState{}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}})
		if !same(out, ctl(toulmin.ControlCaseBegin, 100)) {
			t.Errorf("out = %+v, want CASE-BEGIN 100", out)
		}
		if !st.open || st.idx != 0 || st.tid != 100 {
			t.Errorf("st = %+v, want open idx=0 tid=100", st)
		}
	})

	t.Run("case open, block is the same case: no-op, stays open", func(t *testing.T) {
		st := &caseState{open: true, idx: 0, tid: 100}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}})
		if len(out) != 0 {
			t.Errorf("out = %v, want empty", out)
		}
		if !st.open || st.idx != 0 || st.tid != 100 {
			t.Errorf("st = %+v, want unchanged open idx=0 tid=100", st)
		}
	})

	t.Run("case open, block enters a different case: CASE_END then CASE_BEGIN", func(t *testing.T) {
		st := &caseState{open: true, idx: 0, tid: 100}
		out := run(t, st, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{1}}})
		if !same(out, ctl(toulmin.ControlCaseEnd, 100), ctl(toulmin.ControlCaseBegin, 200)) {
			t.Errorf("out = %+v, want CASE-END 100 then CASE-BEGIN 200", out)
		}
		if !st.open || st.idx != 1 || st.tid != 200 {
			t.Errorf("st = %+v, want open idx=1 tid=200", st)
		}
	})

	failsAtEachID(t, 2, transitionCaseTT(), func(e *enc) error {
		_, err := transitionCase(e, &caseState{open: true, idx: 0, tid: 100}, doc.Block{Kind: doc.StmtBlock, Ref: doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{1}}})
		return err
	})
}

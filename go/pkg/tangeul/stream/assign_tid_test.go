//ff:func feature=tangl type=util control=sequence
//ff:what TestAssignTID — issues ascending TIDs from next and records each under the StmtRef's canonical key
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestAssignTID calls assignTID twice against the same counter and table,
// checking that the returned TID equals the pre-call counter value, the
// counter advances by exactly one, and the table records the TID under the
// ref's tidKey.
func TestAssignTID(t *testing.T) {
	tt := &tidTable{byRef: map[string]toulmin.TID{}}
	next := toulmin.TID(5)

	ref0 := doc.StmtRef{Sec: toulmin.SectionSubject}
	id0 := assignTID(tt, &next, ref0)
	if id0 != 5 {
		t.Errorf("first id = %d, want 5", id0)
	}
	if next != 6 {
		t.Errorf("next after first call = %d, want 6", next)
	}
	if tt.byRef[tidKey(ref0)] != 5 {
		t.Errorf("byRef[%s] = %d, want 5", tidKey(ref0), tt.byRef[tidKey(ref0)])
	}

	ref1 := doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{0}}
	id1 := assignTID(tt, &next, ref1)
	if id1 != 6 {
		t.Errorf("second id = %d, want 6", id1)
	}
	if next != 7 {
		t.Errorf("next after second call = %d, want 7", next)
	}
	if tt.byRef[tidKey(ref1)] != 6 {
		t.Errorf("byRef[%s] = %d, want 6", tidKey(ref1), tt.byRef[tidKey(ref1)])
	}
}

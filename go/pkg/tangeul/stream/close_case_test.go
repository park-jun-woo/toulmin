//ff:func feature=tangl type=util control=sequence
//ff:what TestCloseCase — closeCase is a no-op for a closed case and emits CASE-END for an open one
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// TestCloseCase checks the closed state yields nothing and the open state a
// single CASE-END carrying the open case's Edge Id.
func TestCloseCase(t *testing.T) {
	if got, err := closeCase(testEnc(nil), &caseState{}); err != nil || got != nil {
		t.Fatalf("closed: got %v %v, want nil", got, err)
	}
	got, err := closeCase(testEnc(nil), &caseState{open: true, tid: 7})
	if err != nil {
		t.Fatal(err)
	}
	cps := controlsOf(t, got)
	if len(cps) != 1 || cps[0].Code != toulmin.ControlCaseEnd || cps[0].Arg != 7 {
		t.Fatalf("open: got %+v, want CASE-END 7", cps)
	}
}

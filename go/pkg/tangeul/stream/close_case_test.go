//ff:func feature=tangl type=util control=sequence
//ff:what TestCloseCase — closeCase is a no-op for a closed case and emits CASE-END for an open one
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestCloseCase(t *testing.T) {
	if got := closeCase(testEnc(nil), &caseState{}); got != nil {
		t.Fatalf("closed: got %v, want nil", got)
	}
	cps := controlsOf(t, closeCase(testEnc(nil), &caseState{open: true, tid: 7}))
	if len(cps) != 1 || cps[0].Code != toulmin.ControlCaseEnd || cps[0].Arg != 7 {
		t.Fatalf("open: got %+v, want CASE-END 7", cps)
	}
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestControlWords — controlWords issues a fresh Edge Id and encodes the control and its argument
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// TestControlWords checks a CASE-BEGIN record takes the next Edge Id and carries
// the case REF, and that an argument the width cannot hold is an error.
func TestControlWords(t *testing.T) {
	e := testEnc(nil)
	got, err := controlWords(e, toulmin.ControlCaseBegin, 42, nil)
	if err != nil {
		t.Fatalf("controlWords: %v", err)
	}
	cps := controlsOf(t, got)
	if len(cps) != 1 || cps[0].ID != 1 || cps[0].Code != toulmin.ControlCaseBegin || cps[0].Arg != 42 || e.next != 2 {
		t.Errorf("got %+v next=%d, want one CASE-BEGIN id 1 arg 42", cps, e.next)
	}
	if _, err := controlWords(e, toulmin.ControlCaseEnd, 0x10000, nil); err == nil {
		t.Error("REF past the 1-word width: want error")
	}
}

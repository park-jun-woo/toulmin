//ff:func feature=tangl type=util control=sequence
//ff:what TestControlWords — a CASE-BEGIN with a fresh Edge Id; a bad section code records an error and returns nothing
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestControlWords(t *testing.T) {
	e := testEnc(nil)
	cps := controlsOf(t, controlWords(e, toulmin.ControlCaseBegin, 42, nil))
	if e.err != nil || len(cps) != 1 || cps[0].ID != 1 || cps[0].Code != toulmin.ControlCaseBegin || cps[0].Arg != 42 || e.next != 2 {
		t.Errorf("got %+v next=%d err=%v, want one CASE-BEGIN id 1 arg 42", cps, e.next, e.err)
	}
	if ws := controlWords(e, toulmin.ControlSection, 0x10000, nil); ws != nil || e.err == nil {
		t.Errorf("section code past one word: got %04X err %v", ws, e.err)
	}
}

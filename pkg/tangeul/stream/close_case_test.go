//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCloseCase — closeCase is a no-op when no case is open, else emits a CASE_END matching caseControlWords
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

// TestCloseCase checks both arms of closeCase: a closed caseState yields no
// words and no error, and an open caseState yields exactly the CASE_END words
// caseControlWords would encode for its TID.
func TestCloseCase(t *testing.T) {
	cases := []struct {
		name string
		st   *caseState
	}{
		{"closed", &caseState{open: false}},
		{"open", &caseState{open: true, tid: 7}},
	}
	for _, c := range cases {
		got, err := closeCase(c.st)
		if err != nil {
			t.Fatalf("%s: closeCase: %v", c.name, err)
		}
		if !c.st.open && got != nil {
			t.Errorf("%s: got %v, want nil", c.name, got)
		}
		if !c.st.open {
			continue
		}
		want, err := caseControlWords(codebook.ControlCaseEnd, c.st.tid)
		if err != nil {
			t.Fatalf("%s: caseControlWords: %v", c.name, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v, want %v", c.name, got, want)
		}
	}
}

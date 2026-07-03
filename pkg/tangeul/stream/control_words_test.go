//ff:func feature=tangl type=util control=sequence
//ff:what TestCaseControlWords — caseControlWords round-trips a CASE_BEGIN arg through packet.DecodeControl
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestCaseControlWords encodes a CASE_BEGIN control carrying a Case TID and
// verifies the words decode back to the same code and argument.
func TestCaseControlWords(t *testing.T) {
	got, err := caseControlWords(codebook.ControlCaseBegin, 42)
	if err != nil {
		t.Fatalf("caseControlWords: %v", err)
	}
	cp, n, err := packet.DecodeControl(got)
	if err != nil {
		t.Fatalf("DecodeControl: %v", err)
	}
	if n != len(got) {
		t.Errorf("consumed %d words, want %d", n, len(got))
	}
	if cp.Code != codebook.ControlCaseBegin || cp.Arg != 42 {
		t.Errorf("got Code=%v Arg=%d, want ControlCaseBegin/42", cp.Code, cp.Arg)
	}
}

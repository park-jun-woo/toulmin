//ff:func feature=tangl type=util control=sequence
//ff:what TestQuantOpFromAST — "any" maps to QuantAny, an unknown quantifier string is a hard error
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestQuantOpFromAST exercises both branches of QuantOpFromAST: the sole
// implemented quantifier "any" maps to QuantAny, and any other string is a
// hard error — never a silent fallback.
func TestQuantOpFromAST(t *testing.T) {
	code, err := QuantOpFromAST("any")
	if err != nil || code != toulmin.QuantAny {
		t.Fatalf(`QuantOpFromAST("any"): got 0x%04X err=%v, want 0x%04X`, uint16(code), err, uint16(toulmin.QuantAny))
	}
	if _, err := QuantOpFromAST("all"); err == nil {
		t.Error("QuantOpFromAST(unknown): expected error")
	}
}

//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareOperandError — a malformed left operand (spec without a term) is rejected before the operator
package parser

import (
	"strings"
	"testing"
)

func TestParseCompareOperandError(t *testing.T) {
	_, err := parseCompare("spec equals 3", 9, "t.md")
	if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted spec term after 'spec'") {
		t.Fatalf("err = %v, want the operand error", err)
	}
	if !strings.Contains(err.Error(), "t.md:9") {
		t.Fatalf("err = %v, want the source location", err)
	}
}

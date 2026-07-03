//ff:func feature=tangl type=util control=sequence
//ff:what TestParseCompareValue — parseCompareValue reads NUM and Text slots and rejects truncation, a bad index and an unknown kind
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestParseCompareValue exercises every case of parseCompareValue: a NUM32
// value, a Text label reference, a truncated NUM value, a truncated Text
// value, an out-of-range Text label index, and an unknown ValueKind.
func TestParseCompareValue(t *testing.T) {
	num := word.EncodeNum32(42)
	recs := []word.Word{num[0], num[1]}
	pos := 0
	got, err := parseCompareValue(recs, &pos, nil, ValueNum)
	if err != nil || got != "42" || pos != 2 {
		t.Fatalf("NUM value: got %q, pos %d, err %v", got, pos, err)
	}

	recs = []word.Word{0}
	pos = 0
	got, err = parseCompareValue(recs, &pos, []string{"gold"}, ValueText)
	if err != nil || got != "gold" || pos != 1 {
		t.Fatalf("Text value: got %q, pos %d, err %v", got, pos, err)
	}

	pos = 0
	if _, err := parseCompareValue(nil, &pos, nil, ValueNum); err == nil {
		t.Fatalf("truncated NUM value: want error")
	}

	pos = 0
	if _, err := parseCompareValue(nil, &pos, nil, ValueText); err == nil {
		t.Fatalf("truncated Text value: want error")
	}

	recs = []word.Word{5}
	pos = 0
	if _, err := parseCompareValue(recs, &pos, []string{"only"}, ValueText); err == nil {
		t.Fatalf("out-of-range Text index: want error")
	}

	recs = []word.Word{0}
	pos = 0
	if _, err := parseCompareValue(recs, &pos, nil, ValueKind(9)); err == nil {
		t.Fatalf("unknown ValueKind: want error")
	}
}

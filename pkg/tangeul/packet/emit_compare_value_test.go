//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitCompareValue — exact non-negative uint32 strings inline as NUM32, every other string interns as Text
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEmitCompareValue locks the literal value slot layout: an exact
// non-negative uint32 (identical decimal round-trip) becomes ValueNum plus an
// inline 2-word NUM32; a non-numeric string, a non-canonical spelling ("007"),
// an out-of-range number and a negative number all become ValueText plus an
// interned label index.
func TestEmitCompareValue(t *testing.T) {
	zero := word.EncodeNum32(0)
	num := word.EncodeNum32(42)
	cases := []struct {
		value string
		want  []word.Word
	}{
		{"0", []word.Word{word.Word(ValueNum), zero[0], zero[1]}},
		{"42", []word.Word{word.Word(ValueNum), num[0], num[1]}},
		{"007", []word.Word{word.Word(ValueText), 0}},
		{"4294967296", []word.Word{word.Word(ValueText), 1}},
		{`"gold"`, []word.Word{word.Word(ValueText), 2}},
		{"-1", []word.Word{word.Word(ValueText), 3}},
	}
	var labels []string
	for _, c := range cases {
		got := emitCompareValue(c.value, &labels)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("emitCompareValue(%q) = %v, want %v", c.value, got, c.want)
		}
	}
	if want := []string{"007", "4294967296", `"gold"`, "-1"}; !reflect.DeepEqual(labels, want) {
		t.Fatalf("labels %q, want %q", labels, want)
	}
}

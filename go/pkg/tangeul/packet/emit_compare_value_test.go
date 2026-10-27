//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitCompareValue — exact non-negative uint64 strings as embedded UNSIGNED-INTEGER Values, every other string interns as Text
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestEmitCompareValue locks the literal value slot layout: an exact
// non-negative uint64 (identical decimal round-trip) becomes ValueNum plus an
// embedded UNSIGNED-INTEGER Value in minimal words; a non-numeric string, a non-canonical spelling ("007"),
// an out-of-range number and a negative number all become ValueText plus an
// interned label index.
func TestEmitCompareValue(t *testing.T) {
	zero := uintVal(0)
	num := uintVal(42)
	cases := []struct {
		value string
		want  []core.Word
	}{
		{"0", []core.Word{core.Word(ValueNum), zero[0], zero[1]}},
		{"4294967296", []core.Word{core.Word(ValueNum), 0x0982, 1, 0, 0}},
		{"42", []core.Word{core.Word(ValueNum), num[0], num[1]}},
		{"007", []core.Word{core.Word(ValueText), 0}},
		{"18446744073709551616", []core.Word{core.Word(ValueText), 1}},
		{`"gold"`, []core.Word{core.Word(ValueText), 2}},
		{"-1", []core.Word{core.Word(ValueText), 3}},
	}
	var labels []string
	for _, c := range cases {
		got := emitCompareValue(c.value, &labels)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("emitCompareValue(%q) = %v, want %v", c.value, got, c.want)
		}
	}
	if want := []string{"007", "18446744073709551616", `"gold"`, "-1"}; !reflect.DeepEqual(labels, want) {
		t.Fatalf("labels %q, want %q", labels, want)
	}
}

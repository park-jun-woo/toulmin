//ff:func feature=tangl type=util control=sequence
//ff:what TestStructField — tests exact, case-insensitive, unexported, and missing field lookup
package eval

import (
	"reflect"
	"testing"
)

func TestStructField(t *testing.T) {
	v := reflect.ValueOf(struct {
		TrustScore float64
		hidden     int
	}{TrustScore: 0.8, hidden: 1})

	if fv, ok := structField(v, "TrustScore"); !ok || fv.Float() != 0.8 {
		t.Fatalf("exact match: got %v,%v", fv, ok)
	}
	if fv, ok := structField(v, "trustScore"); !ok || fv.Float() != 0.8 {
		t.Fatalf("case-insensitive match: got %v,%v", fv, ok)
	}
	if _, ok := structField(v, "hidden"); ok {
		t.Fatal("unexported field should not resolve")
	}
	if _, ok := structField(v, "ghost"); ok {
		t.Fatal("missing field should not resolve")
	}
}

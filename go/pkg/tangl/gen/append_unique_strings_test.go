//ff:func feature=tangl type=codegen control=sequence
//ff:what TestAppendUniqueStrings — values not already in dst (or earlier in values) are appended in order
package gen

import (
	"reflect"
	"testing"
)

func TestAppendUniqueStrings(t *testing.T) {
	got := appendUniqueStrings([]string{"a", "b"}, "b", "c", "a", "c", "d")
	if want := []string{"a", "b", "c", "d"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("appendUniqueStrings = %v, want %v", got, want)
	}
	if got := appendUniqueStrings(nil); got != nil {
		t.Fatalf("nothing to add = %v, want nil", got)
	}
	if got := appendUniqueStrings(nil, "x", "x"); !reflect.DeepEqual(got, []string{"x"}) {
		t.Fatalf("duplicate values = %v", got)
	}
}

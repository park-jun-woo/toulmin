//ff:func feature=tangl type=engine control=sequence
//ff:what TestUniqueTerms — groups concatenate in order keeping first occurrences
package eval

import (
	"reflect"
	"testing"
)

func TestUniqueTerms(t *testing.T) {
	got := uniqueTerms([]string{"a", "b"}, nil, []string{"b", "c", "a"})
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("uniqueTerms = %v, want %v", got, want)
	}
	if got := uniqueTerms(); got != nil {
		t.Fatalf("no groups = %v, want nil", got)
	}
}

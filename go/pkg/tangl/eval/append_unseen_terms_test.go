//ff:func feature=tangl type=engine control=sequence
//ff:what TestAppendUnseenTerms — unseen terms append and are marked; seen ones are skipped
package eval

import (
	"reflect"
	"testing"
)

func TestAppendUnseenTerms(t *testing.T) {
	seen := map[string]bool{"a": true}
	got := appendUnseenTerms([]string{"a"}, seen, []string{"a", "b", "b", "c"})
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("appendUnseenTerms = %v, want %v", got, want)
	}
	if !seen["b"] || !seen["c"] {
		t.Fatalf("seen = %v, want b and c marked", seen)
	}
}

//ff:func feature=tangl type=util control=sequence
//ff:what TestEqualSliceIgnoringLine — slice length guard, nil/empty equality, per-element mismatch and full match
package tangeul

import (
	"reflect"
	"testing"
)

// TestEqualSliceIgnoringLine drives every branch of equalSliceIgnoringLine
// directly: a length mismatch (the leading false return before the loop), a nil
// slice against an empty slice (both length 0, so the loop body never runs and
// the two representationally-different-but-semantically-equal slices compare
// equal), a slice pair whose first element differs (the early false return
// mid-loop), and a slice pair whose every element is equal (the loop runs to
// completion and returns true).
func TestEqualSliceIgnoringLine(t *testing.T) {
	t.Run("LengthMismatch", func(t *testing.T) {
		x := []int{1, 2}
		y := []int{1}
		if equalSliceIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected false: differing lengths")
		}
	})

	t.Run("NilEqualsEmpty", func(t *testing.T) {
		var x []int
		y := []int{}
		if !equalSliceIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected true: nil and empty slices both have length 0")
		}
	})

	t.Run("ElementMismatchReturnsFalseEarly", func(t *testing.T) {
		x := []int{1, 2, 3}
		y := []int{1, 9, 3}
		if equalSliceIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected false: element 1 differs")
		}
	})

	t.Run("AllElementsEqual", func(t *testing.T) {
		x := []int{1, 2, 3}
		y := []int{1, 2, 3}
		if !equalSliceIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected true: every element equal")
		}
	})
}

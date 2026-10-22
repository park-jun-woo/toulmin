//ff:func feature=tangl type=util control=sequence
//ff:what TestEqualIgnoringLineValue — equalIgnoringLine's type-mismatch guard and its Struct/Slice/default dispatch
package tangeul

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEqualIgnoringLineValue drives the non-nil-handling branches of
// equalIgnoringLine: the leading type-mismatch guard, the Struct and Slice
// dispatch to their own helpers, and the default reflect.DeepEqual fallback for
// scalar and map kinds. The Pointer/Interface nil-combination branches share the
// same case body and are covered separately in
// TestEqualIgnoringLineValueNilCombinations.
func TestEqualIgnoringLineValue(t *testing.T) {
	t.Run("TypeMismatch", func(t *testing.T) {
		if equalIgnoringLine(reflect.ValueOf(1), reflect.ValueOf("x")) {
			t.Fatal("expected false: int and string are different types")
		}
	})

	t.Run("StructDispatch", func(t *testing.T) {
		x := ast.See{Alias: "bank", Package: "bank/core"}
		y := ast.See{Alias: "bank", Package: "bank/core"}
		if !equalIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected true: equal See structs delegate to equalStructIgnoringLine")
		}
		y.Package = "bank/other"
		if equalIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected false: unequal See structs delegate to equalStructIgnoringLine")
		}
	})

	t.Run("SliceDispatch", func(t *testing.T) {
		x := []ast.See{{Alias: "a"}, {Alias: "b"}}
		y := []ast.See{{Alias: "a"}, {Alias: "b"}}
		if !equalIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected true: equal slices delegate to equalSliceIgnoringLine")
		}
		y[1].Alias = "c"
		if equalIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected false: unequal slices delegate to equalSliceIgnoringLine")
		}
	})

	t.Run("DefaultDeepEqual", func(t *testing.T) {
		if !equalIgnoringLine(reflect.ValueOf(5), reflect.ValueOf(5)) {
			t.Fatal("expected true: equal ints via DeepEqual fallback")
		}
		if equalIgnoringLine(reflect.ValueOf(5), reflect.ValueOf(6)) {
			t.Fatal("expected false: unequal ints via DeepEqual fallback")
		}
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		if !equalIgnoringLine(reflect.ValueOf(m1), reflect.ValueOf(m2)) {
			t.Fatal("expected true: equal maps via DeepEqual fallback")
		}
	})
}

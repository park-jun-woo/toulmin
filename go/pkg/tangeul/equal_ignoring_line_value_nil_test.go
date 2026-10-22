//ff:func feature=tangl type=util control=sequence
//ff:what TestEqualIgnoringLineValueNilCombinations — equalIgnoringLine's shared Pointer/Interface nil handling
package tangeul

import (
	"reflect"
	"testing"
)

// TestEqualIgnoringLineValueNilCombinations drives the nil-handling half of
// equalIgnoringLine's shared Pointer/Interface case: both nil, exactly one nil
// in either order, and neither nil with both an equal and an unequal recursive
// outcome. Pointer kind is exercised with reflect.ValueOf on *int values
// directly; Interface kind needs a struct field typed interface{} because a
// bare reflect.ValueOf(x) always unwraps to x's concrete type.
func TestEqualIgnoringLineValueNilCombinations(t *testing.T) {
	t.Run("PointerBothNil", func(t *testing.T) {
		var a, b *int
		if !equalIgnoringLine(reflect.ValueOf(a), reflect.ValueOf(b)) {
			t.Fatal("expected true: both nil pointers")
		}
	})

	t.Run("PointerOneNil", func(t *testing.T) {
		n := 1
		var nilP *int
		if equalIgnoringLine(reflect.ValueOf(&n), reflect.ValueOf(nilP)) {
			t.Fatal("expected false: non-nil x, nil y")
		}
		if equalIgnoringLine(reflect.ValueOf(nilP), reflect.ValueOf(&n)) {
			t.Fatal("expected false: nil x, non-nil y")
		}
	})

	t.Run("PointerBothNonNil", func(t *testing.T) {
		a, b, c := 7, 7, 8
		if !equalIgnoringLine(reflect.ValueOf(&a), reflect.ValueOf(&b)) {
			t.Fatal("expected true: pointers to equal ints recurse to equal")
		}
		if equalIgnoringLine(reflect.ValueOf(&a), reflect.ValueOf(&c)) {
			t.Fatal("expected false: pointers to unequal ints recurse to unequal")
		}
	})

	t.Run("InterfaceNilCombinations", func(t *testing.T) {
		type ifaceHolder struct{ V interface{} }

		bothNilX := reflect.ValueOf(ifaceHolder{}).Field(0)
		bothNilY := reflect.ValueOf(ifaceHolder{}).Field(0)
		if bothNilX.Kind() != reflect.Interface {
			t.Fatalf("fixture broken: Kind = %v, want Interface", bothNilX.Kind())
		}
		if !equalIgnoringLine(bothNilX, bothNilY) {
			t.Fatal("expected true: both nil interfaces")
		}

		set := reflect.ValueOf(ifaceHolder{V: 1}).Field(0)
		unset := reflect.ValueOf(ifaceHolder{}).Field(0)
		if equalIgnoringLine(set, unset) {
			t.Fatal("expected false: non-nil x, nil y interface")
		}
		if equalIgnoringLine(unset, set) {
			t.Fatal("expected false: nil x, non-nil y interface")
		}

		equalOther := reflect.ValueOf(ifaceHolder{V: 1}).Field(0)
		if !equalIgnoringLine(set, equalOther) {
			t.Fatal("expected true: equal non-nil interface values recurse to equal")
		}
		diffOther := reflect.ValueOf(ifaceHolder{V: 2}).Field(0)
		if equalIgnoringLine(set, diffOther) {
			t.Fatal("expected false: unequal non-nil interface values recurse to unequal")
		}
	})
}

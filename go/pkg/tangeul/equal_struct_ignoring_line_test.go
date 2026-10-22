//ff:func feature=tangl type=util control=sequence
//ff:what TestEqualStructIgnoringLine — struct field loop: zero fields, Line skip, per-field mismatch and full match
package tangeul

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEqualStructIgnoringLine drives every branch of equalStructIgnoringLine
// directly: a zero-field struct (the loop body never runs), a struct whose only
// difference is its "Line" field (the continue-skip branch, reaching the final
// true with the skip taken), a struct whose first non-Line field differs (the
// early false return before the loop finishes), and a struct whose every field
// -- Line included -- is equal (the loop runs to completion and returns true).
// ast.See supplies a real three-field struct (Alias, Package, Line) so the
// fixture matches the production shape the function actually walks.
func TestEqualStructIgnoringLine(t *testing.T) {
	t.Run("ZeroFields", func(t *testing.T) {
		type empty struct{}
		if !equalStructIgnoringLine(reflect.ValueOf(empty{}), reflect.ValueOf(empty{})) {
			t.Fatal("expected true: a zero-field struct has nothing to disagree on")
		}
	})

	t.Run("SkipsLineFieldDifference", func(t *testing.T) {
		x := ast.See{Alias: "bank", Package: "bank/core", Line: 1}
		y := ast.See{Alias: "bank", Package: "bank/core", Line: 99}
		if !equalStructIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected true: only the skipped Line field differs")
		}
	})

	t.Run("FieldMismatchReturnsFalseEarly", func(t *testing.T) {
		x := ast.See{Alias: "bank", Package: "bank/core", Line: 1}
		y := ast.See{Alias: "other", Package: "bank/core", Line: 1}
		if equalStructIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected false: Alias field differs")
		}
	})

	t.Run("AllFieldsEqual", func(t *testing.T) {
		x := ast.See{Alias: "bank", Package: "bank/core", Line: 5}
		y := ast.See{Alias: "bank", Package: "bank/core", Line: 5}
		if !equalStructIgnoringLine(reflect.ValueOf(x), reflect.ValueOf(y)) {
			t.Fatal("expected true: every field, including Line, is equal")
		}
	})
}

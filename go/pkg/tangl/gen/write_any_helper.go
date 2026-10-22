//ff:func feature=tangl type=codegen control=sequence
//ff:what writeAnyHelper — writes the shared tanglAny function
package gen

import "strings"

// writeAnyHelper writes the shared tanglAny function: the exists quantifier
// runtime that iterates a slice/array source in index order and returns true
// with the satisfying element's witness at the first hit, mirroring
// pkg/tangl/eval.compileExists. A no-value or non-slice source counts as zero
// elements; when no element satisfies, the witness is a Fact whose Actual is
// the number of elements examined.
func writeAnyHelper(w *strings.Builder) {
	w.WriteString(`func tanglAny(src any, ok bool, where, expected string, elem func(any, int) (bool, any)) (bool, any) {
	if !ok {
		return false, toulmin.Fact{Where: where, Expected: expected, Actual: 0}
	}
	rv := reflect.ValueOf(src)
	if k := rv.Kind(); k != reflect.Slice && k != reflect.Array {
		return false, toulmin.Fact{Where: where, Expected: expected, Actual: 0}
	}
	n := rv.Len()
	for i := 0; i < n; i++ {
		if b, w := elem(rv.Index(i).Interface(), i); b {
			return true, w
		}
	}
	return false, toulmin.Fact{Where: where, Expected: expected, Actual: n}
}

`)
}

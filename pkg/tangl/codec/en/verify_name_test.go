//ff:func feature=tangl type=parser control=sequence
//ff:what TestVerifyName — a matching slot passes, an out-of-range or mismatched name fails
package en

import "testing"

// TestVerifyName covers all three outcomes: a name equal to slice[i] passes, an
// index at or beyond len(slice) fails (the `i >= len(slice)` half of the guard),
// and an in-range name that disagrees with slice[i] fails (the `slice[i] != name`
// half of the guard).
func TestVerifyName(t *testing.T) {
	slice := []string{"a", "b"}
	if err := verifyName(slice, 0, "a"); err != nil {
		t.Errorf("verifyName(slice, 0, \"a\") = %v, want nil", err)
	}
	if err := verifyName(slice, 2, "a"); err == nil {
		t.Error("verifyName accepted an out-of-range index")
	}
	if err := verifyName(slice, 1, "wrong"); err == nil {
		t.Error("verifyName accepted a name that disagrees with the slice")
	}
}

//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestNumericOp — tests every numeric operator in both outcomes
package eval

import "testing"

func TestNumericOp(t *testing.T) {
	cases := []struct {
		op     string
		fv, tv float64
		want   bool
	}{
		{"is greater than", 4, 3, true},
		{"is greater than", 3, 3, false},
		{"is less than", 2, 3, true},
		{"is less than", 3, 3, false},
		{"is at most", 3, 3, true},
		{"is at most", 4, 3, false},
		{"is at least", 3, 3, true},
		{"is at least", 2, 3, false},
	}
	for _, c := range cases {
		if got := numericOp(c.op, c.fv, c.tv); got != c.want {
			t.Fatalf("numericOp(%q,%v,%v)=%v want %v", c.op, c.fv, c.tv, got, c.want)
		}
	}
}

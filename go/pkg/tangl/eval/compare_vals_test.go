//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestCompareVals — tests every tangl comparison operator over resolved values
package eval

import "testing"

func TestCompareVals(t *testing.T) {
	cases := []struct {
		name string
		lhs  any
		ok   bool
		op   string
		rhs  any
		want bool
	}{
		{"IsEmptyNoValue", nil, false, "is empty", "", true},
		{"IsEmptyEmptyString", "", true, "is empty", "", true},
		{"IsEmptyFalse", "x", true, "is empty", "", false},
		{"IsNotEmptyTrue", "x", true, "is not empty", "", true},
		{"IsNotEmptyNoValue", nil, false, "is not empty", "", false},
		{"EqualsTrue", "abc", true, "equals", "abc", true},
		{"EqualsNumericStringify", 5, true, "equals", "5", true},
		{"EqualsNoValue", nil, false, "equals", "abc", false},
		{"IsInMatch", "b", true, "is in", "a, b, c", true},
		{"IsInNoValue", nil, false, "is in", "a, b, c", false},
		{"IsInNoMatch", "z", true, "is in", "a, b, c", false},
		{"GreaterTrue", 10.0, true, "is greater than", 3, true},
		{"GreaterFalse", 2.0, true, "is greater than", 3, false},
		{"LessTrue", 2.0, true, "is less than", 3, true},
		{"AtMostTrue", 3.0, true, "is at most", 3, true},
		{"AtLeastTrue", 3.0, true, "is at least", 3, true},
		{"NumericNoValue", nil, false, "is greater than", 3, false},
		{"NumericUnparseable", "abc", true, "is greater than", "3", false},
		{"UnknownOp", "x", true, "is whatever", "x", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := compareVals(c.lhs, c.ok, c.op, c.rhs); got != c.want {
				t.Fatalf("compareVals(%v,%v,%q,%v)=%v want %v", c.lhs, c.ok, c.op, c.rhs, got, c.want)
			}
		})
	}
}

//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompareNumeric — tests float parsing of both operands and unparseable falls
package eval

import "testing"

func TestCompareNumeric(t *testing.T) {
	if !compareNumeric("is greater than", 10, "3") {
		t.Fatal("10 > 3 should hold across stringified forms")
	}
	if compareNumeric("is greater than", 2, 3) {
		t.Fatal("2 > 3 should not hold")
	}
	if compareNumeric("is greater than", "abc", 3) {
		t.Fatal("unparseable lhs should be false")
	}
	if compareNumeric("is greater than", 10, "abc") {
		t.Fatal("unparseable rhs should be false")
	}
}

//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompareIn — tests comma-list membership with trimming and a miss
package eval

import "testing"

func TestCompareIn(t *testing.T) {
	if !compareIn("b", "a, b, c") {
		t.Fatal("b should be in the trimmed comma list")
	}
	if compareIn("z", "a, b, c") {
		t.Fatal("z should not be in the list")
	}
}

//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCompareOpFromAST — every surface value maps to a code that maps back; an unknown value fails
package codebook

import (
	"testing"
)

func TestCompareOpFromAST(t *testing.T) {
	for _, v := range []string{"is empty", "is not empty", "equals", "is in", "is greater than", "is less than", "is at most", "is at least"} {
		c, err := CompareOpFromAST(v)
		if err != nil {
			t.Fatalf("%v: %v", v, err)
		}
		if back, err := CompareOpToAST(c); err != nil || back != v {
			t.Fatalf("%v: round trip got %v %v", v, back, err)
		}
	}
	if _, err := CompareOpFromAST("bogus"); err == nil {
		t.Fatal("unknown value: want error")
	}
}

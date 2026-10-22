//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCompareOpToAST — every code maps to its surface value; an unknown code fails
package codebook

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestCompareOpToAST(t *testing.T) {
	for _, v := range []string{"is empty", "is not empty", "equals", "is in", "is greater than", "is less than", "is at most", "is at least"} {
		c, _ := CompareOpFromAST(v)
		if got, err := CompareOpToAST(c); err != nil || got != v {
			t.Fatalf("%v: got %v %v", v, got, err)
		}
	}
	if _, err := CompareOpToAST(toulmin.CompareOp(0x9999)); err == nil {
		t.Fatal("unknown code: want error")
	}
}

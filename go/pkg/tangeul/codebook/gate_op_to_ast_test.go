//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGateOpToAST — every code maps to its surface value; an unknown code fails
package codebook

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestGateOpToAST(t *testing.T) {
	for _, v := range []string{"at least", "above", "less than", "at most"} {
		c, _ := GateOpFromAST(v)
		if got, err := GateOpToAST(c); err != nil || got != v {
			t.Fatalf("%v: got %v %v", v, got, err)
		}
	}
	if _, err := GateOpToAST(toulmin.GateOp(9)); err == nil {
		t.Fatal("unknown code: want error")
	}
}

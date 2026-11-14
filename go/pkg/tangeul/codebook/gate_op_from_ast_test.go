//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGateOpFromAST — every surface value maps to a code that maps back; an unknown value fails
package codebook

import (
	"testing"
)

func TestGateOpFromAST(t *testing.T) {
	for _, v := range []string{"at least", "above", "less than", "at most"} {
		c, err := GateOpFromAST(v)
		if err != nil {
			t.Fatalf("%v: %v", v, err)
		}
		if back, err := GateOpToAST(c); err != nil || back != v {
			t.Fatalf("%v: round trip got %v %v", v, back, err)
		}
	}
	if _, err := GateOpFromAST("bogus"); err == nil {
		t.Fatal("unknown value: want error")
	}
}

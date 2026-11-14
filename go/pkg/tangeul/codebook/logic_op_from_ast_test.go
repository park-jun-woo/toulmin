//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLogicOpFromAST — every surface value maps to a code that maps back; an unknown value fails
package codebook

import (
	"testing"
)

func TestLogicOpFromAST(t *testing.T) {
	for _, v := range []string{"and", "or"} {
		c, err := LogicOpFromAST(v)
		if err != nil {
			t.Fatalf("%v: %v", v, err)
		}
		if back, err := LogicOpToAST(c); err != nil || back != v {
			t.Fatalf("%v: round trip got %v %v", v, back, err)
		}
	}
	if _, err := LogicOpFromAST("xor"); err == nil {
		t.Fatal("unknown value: want error")
	}
}

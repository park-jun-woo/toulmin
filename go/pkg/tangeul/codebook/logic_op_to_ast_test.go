//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLogicOpToAST — every code maps to its surface value; an unknown code fails
package codebook

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestLogicOpToAST(t *testing.T) {
	for _, v := range []string{"and", "or"} {
		c, _ := LogicOpFromAST(v)
		if got, err := LogicOpToAST(c); err != nil || got != v {
			t.Fatalf("%v: got %v %v", v, got, err)
		}
	}
	if _, err := LogicOpToAST(toulmin.LogicOp(0x9999)); err == nil {
		t.Fatal("unknown code: want error")
	}
}

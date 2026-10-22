//ff:func feature=tangl type=parser control=sequence
//ff:what TestCompareLhsTextKo — bare `field` for a nil Lhs vs rendered operand for a set Lhs
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCompareLhsTextKo(t *testing.T) {
	t.Run("NilLhsBareField", func(t *testing.T) {
		text, base, err := compareLhsTextKo(ast.Compare{Field: "금액"})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if text != "`금액`" || base != "금액" {
			t.Fatalf("got text=%q base=%q", text, base)
		}
	})

	t.Run("LhsOperand", func(t *testing.T) {
		c := ast.Compare{Lhs: ast.Path{Steps: []string{"c", "최소가격"}}}
		text, base, err := compareLhsTextKo(c)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if text != "`c` 의 `최소가격`" || base != "최소가격" {
			t.Fatalf("got text=%q base=%q", text, base)
		}
	})
}

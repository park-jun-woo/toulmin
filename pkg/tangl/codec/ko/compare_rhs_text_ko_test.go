//ff:func feature=tangl type=parser control=sequence
//ff:what TestCompareRhsTextKo — verbatim Value for a nil Rhs vs rendered operand for a set Rhs
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCompareRhsTextKo(t *testing.T) {
	t.Run("NilRhsVerbatimValue", func(t *testing.T) {
		text, base, err := compareRhsTextKo(ast.Compare{Value: "650"})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if text != "650" || base != "650" {
			t.Fatalf("got text=%q base=%q", text, base)
		}
	})

	t.Run("RhsOperand", func(t *testing.T) {
		c := ast.Compare{Rhs: ast.Path{Steps: []string{"기본가"}}}
		text, base, err := compareRhsTextKo(c)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if text != "`기본가`" || base != "기본가" {
			t.Fatalf("got text=%q base=%q", text, base)
		}
	})
}

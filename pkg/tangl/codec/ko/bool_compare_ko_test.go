//ff:func feature=tangl type=parser control=sequence
//ff:what TestBoolCompareKo — single-step demotion to Field vs multi-step Lhs path
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBoolCompareKo(t *testing.T) {
	t.Run("SingleStepToField", func(t *testing.T) {
		cmp := boolCompareKo([]string{"활성"})
		if cmp.Field != "활성" || cmp.Op != "equals" || cmp.Value != "true" {
			t.Fatalf("got %#v", cmp)
		}
		if cmp.Lhs != nil {
			t.Fatalf("Lhs must stay nil for a single-step path, got %#v", cmp.Lhs)
		}
	})

	t.Run("MultiStepToLhs", func(t *testing.T) {
		cmp := boolCompareKo([]string{"작성자", "활성"})
		if cmp.Field != "" || cmp.Op != "equals" || cmp.Value != "true" {
			t.Fatalf("got %#v", cmp)
		}
		p, ok := cmp.Lhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || p.Steps[0] != "작성자" || p.Steps[1] != "활성" {
			t.Fatalf("Lhs = %#v", cmp.Lhs)
		}
	})
}

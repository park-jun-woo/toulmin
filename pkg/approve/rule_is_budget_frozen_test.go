//ff:func feature=approve type=rule control=sequence
//ff:what TestIsBudgetFrozenEdgeCases — IsBudgetFrozen covers the non-*Budget context value branch
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsBudgetFrozenEdgeCases covers the branch of IsBudgetFrozen not
// exercised by TestIsBudgetFrozen's frozen/not-frozen cases: a missing or
// non-*Budget "budget" context value fails the type assertion and returns
// false.
func TestIsBudgetFrozenEdgeCases(t *testing.T) {
	t.Run("wrong type returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("budget", "not-a-budget")
		got, val := IsBudgetFrozen(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing budget returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		got, val := IsBudgetFrozen(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})
}

//ff:func feature=approve type=rule control=sequence
//ff:what TestIsUnderBudgetEdgeCases — IsUnderBudget covers the non-float64 amount and non-*Budget budget branches
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsUnderBudgetEdgeCases covers the two branches of IsUnderBudget not
// exercised by TestIsUnderBudget's under/equal/over cases: a non-float64
// (or missing) "amount" context value fails its type assertion, and a
// non-*Budget (or missing) "budget" context value fails its type assertion;
// both return false.
func TestIsUnderBudgetEdgeCases(t *testing.T) {
	t.Run("non-float64 amount returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("amount", "5000")
		ctx.Set("budget", &Budget{Remaining: 10000})
		got, val := IsUnderBudget(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing amount returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("budget", &Budget{Remaining: 10000})
		got, _ := IsUnderBudget(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("non-*Budget budget returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("amount", 5000.0)
		ctx.Set("budget", "not-a-budget")
		got, val := IsUnderBudget(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing budget returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("amount", 5000.0)
		got, _ := IsUnderBudget(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
	})
}

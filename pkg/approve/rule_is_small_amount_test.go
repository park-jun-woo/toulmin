//ff:func feature=approve type=rule control=sequence
//ff:what TestIsSmallAmountEdgeCases — IsSmallAmount covers the empty-specs and non-float64 amount branches
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsSmallAmountEdgeCases covers the two branches of IsSmallAmount not
// exercised by TestIsSmallAmount's small/equal/large cases: an empty specs
// slice short-circuits to false, and a non-float64 (or missing) "amount"
// context value fails the type assertion and returns false.
func TestIsSmallAmountEdgeCases(t *testing.T) {
	t.Run("empty specs returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("amount", 5000.0)
		got, val := IsSmallAmount(ctx, toulmin.Specs{})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("non-float64 amount returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("amount", "5000")
		got, val := IsSmallAmount(ctx, toulmin.Specs{&ThresholdSpec{Max: 10000}})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing amount returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		got, val := IsSmallAmount(ctx, toulmin.Specs{&ThresholdSpec{Max: 10000}})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})
}

//ff:func feature=approve type=rule control=sequence
//ff:what TestIsCEOOverrideEdgeCases — IsCEOOverride covers the non-string approverRole branch
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsCEOOverrideEdgeCases covers the branch of IsCEOOverride not
// exercised by TestIsCEOOverride's ceo/not-ceo cases: a missing or
// non-string "approverRole" context value fails the type assertion and
// returns false.
func TestIsCEOOverrideEdgeCases(t *testing.T) {
	t.Run("wrong type returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("approverRole", 42)
		got, val := IsCEOOverride(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing approverRole returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		got, val := IsCEOOverride(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})
}

//ff:func feature=approve type=rule control=sequence
//ff:what TestIsAboveLevelEdgeCases — IsAboveLevel covers the empty-specs and non-int approverLevel branches
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsAboveLevelEdgeCases covers the two branches of IsAboveLevel not
// exercised by TestIsAboveLevel's above/equal/below cases: an empty specs
// slice short-circuits to false, and a non-int (or missing) "approverLevel"
// context value fails the type assertion and returns false.
func TestIsAboveLevelEdgeCases(t *testing.T) {
	t.Run("empty specs returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("approverLevel", 5)
		got, val := IsAboveLevel(ctx, toulmin.Specs{})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("non-int approverLevel returns false", func(t *testing.T) {
		ab := &ApproverSpec{Level: 3}
		ctx := toulmin.NewContext()
		ctx.Set("approverLevel", "high")
		got, val := IsAboveLevel(ctx, toulmin.Specs{ab})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("missing approverLevel returns false", func(t *testing.T) {
		ab := &ApproverSpec{Level: 3}
		ctx := toulmin.NewContext()
		got, val := IsAboveLevel(ctx, toulmin.Specs{ab})
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})
}

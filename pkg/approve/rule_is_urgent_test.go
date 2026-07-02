//ff:func feature=approve type=rule control=sequence
//ff:what TestIsUrgentEdgeCases — IsUrgent covers the non-map requestMetadata branch
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestIsUrgentEdgeCases covers the branch of IsUrgent not exercised by
// TestIsUrgent's urgent/not-urgent/nil-meta cases: a non-map (or entirely
// unset) "requestMetadata" context value fails the type assertion and
// returns false.
func TestIsUrgentEdgeCases(t *testing.T) {
	t.Run("wrong type returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("requestMetadata", "not-a-map")
		got, val := IsUrgent(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("unset requestMetadata returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		got, val := IsUrgent(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})

	t.Run("urgent key with non-bool value returns false", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("requestMetadata", map[string]any{"urgent": "yes"})
		got, val := IsUrgent(ctx, nil)
		if got {
			t.Errorf("got %v, want false", got)
		}
		if val != nil {
			t.Errorf("val = %v, want nil", val)
		}
	})
}

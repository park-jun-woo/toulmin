//ff:func feature=price type=rule control=sequence
//ff:what TestIsAlreadyDiscounted_Branches — covers requestMetadata type-assertion failure branch of IsAlreadyDiscounted
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsAlreadyDiscounted_Branches(t *testing.T) {
	t.Run("metadata wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("requestMetadata", "not-a-map")

		got, evidence := IsAlreadyDiscounted(ctx, nil)
		if got {
			t.Errorf("expected false when requestMetadata is not a map[string]any, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("metadata unset", func(t *testing.T) {
		ctx := toulmin.NewContext()

		got, _ := IsAlreadyDiscounted(ctx, nil)
		if got {
			t.Errorf("expected false when requestMetadata is unset, got %v", got)
		}
	})

	t.Run("discounted wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("requestMetadata", map[string]any{"discounted": "yes"})

		got, _ := IsAlreadyDiscounted(ctx, nil)
		if got {
			t.Errorf("expected false when discounted field is not a bool, got %v", got)
		}
	})
}

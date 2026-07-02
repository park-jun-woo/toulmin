//ff:func feature=price type=rule control=sequence
//ff:what TestIsBulkOrder_Branches — covers empty specs and quantity type-assertion failure branches of IsBulkOrder
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsBulkOrder_Branches(t *testing.T) {
	t.Run("empty specs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("quantity", 100)

		got, evidence := IsBulkOrder(ctx, toulmin.Specs{})
		if got {
			t.Errorf("expected false for empty specs, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("quantity wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("quantity", "not-an-int")

		got, _ := IsBulkOrder(ctx, toulmin.Specs{&BulkOrderSpec{MinQuantity: 50}})
		if got {
			t.Errorf("expected false when quantity is not an int, got %v", got)
		}
	})

	t.Run("quantity unset", func(t *testing.T) {
		ctx := toulmin.NewContext()

		got, _ := IsBulkOrder(ctx, toulmin.Specs{&BulkOrderSpec{MinQuantity: 50}})
		if got {
			t.Errorf("expected false when quantity is unset, got %v", got)
		}
	})
}

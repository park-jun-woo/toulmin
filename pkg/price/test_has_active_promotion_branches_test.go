//ff:func feature=price type=rule control=sequence
//ff:what TestHasActivePromotion_Branches — covers empty specs and promotions type-assertion failure branches of HasActivePromotion
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasActivePromotion_Branches(t *testing.T) {
	t.Run("empty specs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("promotions", []Promotion{{Name: "blackfriday", Active: true}})

		got, evidence := HasActivePromotion(ctx, toulmin.Specs{})
		if got {
			t.Errorf("expected false for empty specs, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("promotions wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("promotions", "not-a-slice")

		db := &DiscountSpec{Name: "blackfriday", Fixed: 5000}
		got, _ := HasActivePromotion(ctx, toulmin.Specs{db})
		if got {
			t.Errorf("expected false when promotions is not a []Promotion, got %v", got)
		}
	})

	t.Run("promotions unset", func(t *testing.T) {
		ctx := toulmin.NewContext()

		db := &DiscountSpec{Name: "blackfriday", Fixed: 5000}
		got, _ := HasActivePromotion(ctx, toulmin.Specs{db})
		if got {
			t.Errorf("expected false when promotions is unset, got %v", got)
		}
	})
}

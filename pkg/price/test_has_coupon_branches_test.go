//ff:func feature=price type=rule control=sequence
//ff:what TestHasCoupon_Branches — covers empty specs, basePrice type-assertion failure, and coupons type-assertion failure branches of HasCoupon
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasCoupon_Branches(t *testing.T) {
	t.Run("empty specs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("basePrice", 10000.0)
		ctx.Set("coupons", []Coupon{{Code: "SAVE30", MinPrice: 5000}})

		got, evidence := HasCoupon(ctx, toulmin.Specs{})
		if got {
			t.Errorf("expected false for empty specs, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("base price wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("basePrice", "not-a-float")
		ctx.Set("coupons", []Coupon{{Code: "SAVE30", MinPrice: 5000}})

		db := &DiscountSpec{Name: "SAVE30", Rate: 0.3}
		got, _ := HasCoupon(ctx, toulmin.Specs{db})
		if got {
			t.Errorf("expected false when basePrice is not a float64, got %v", got)
		}
	})

	t.Run("coupons wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("basePrice", 10000.0)
		ctx.Set("coupons", "not-a-slice")

		db := &DiscountSpec{Name: "SAVE30", Rate: 0.3}
		got, _ := HasCoupon(ctx, toulmin.Specs{db})
		if got {
			t.Errorf("expected false when coupons is not a []Coupon, got %v", got)
		}
	})

	t.Run("base price unset", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("coupons", []Coupon{{Code: "SAVE30", MinPrice: 5000}})

		db := &DiscountSpec{Name: "SAVE30", Rate: 0.3}
		got, _ := HasCoupon(ctx, toulmin.Specs{db})
		if got {
			t.Errorf("expected false when basePrice is unset, got %v", got)
		}
	})
}

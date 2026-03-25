//ff:func feature=price type=rule control=iteration dimension=1
//ff:what HasCoupon: spec(DiscountSpec)의 쿠폰 적용 조건 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasCoupon checks if a coupon applies. spec is *DiscountSpec.
func HasCoupon(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	basePrice, _ := ctx.Get("basePrice")
	coupons, _ := ctx.Get("coupons")
	db := specs[0].(*DiscountSpec)
	for _, c := range coupons.([]Coupon) {
		if basePrice.(float64) >= c.MinPrice {
			return true, db
		}
	}
	return false, nil
}

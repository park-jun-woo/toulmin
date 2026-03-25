//ff:func feature=price type=rule control=iteration dimension=1
//ff:what HasCoupon: backing(DiscountBacking)의 쿠폰 적용 조건 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasCoupon checks if a coupon applies. backing is *DiscountBacking.
func HasCoupon(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	basePrice, _ := ctx.Get("basePrice")
	coupons, _ := ctx.Get("coupons")
	db := backing.(*DiscountBacking)
	for _, c := range coupons.([]Coupon) {
		if basePrice.(float64) >= c.MinPrice {
			return true, db
		}
	}
	return false, nil
}

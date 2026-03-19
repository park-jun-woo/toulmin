//ff:func feature=price type=rule control=iteration dimension=1
//ff:what HasCoupon: backing(DiscountBacking)의 쿠폰 적용 조건 판정
package price

// HasCoupon checks if a coupon applies. backing is *DiscountBacking.
func HasCoupon(claim any, ground any, backing any) (bool, any) {
	req := claim.(*PurchaseRequest)
	ctx := ground.(*PriceContext)
	db := backing.(*DiscountBacking)
	for _, c := range ctx.Coupons {
		if req.BasePrice >= c.MinPrice {
			return true, db
		}
	}
	return false, nil
}

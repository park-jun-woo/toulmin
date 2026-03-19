//ff:func feature=price type=rule control=iteration dimension=1
//ff:what HasActivePromotion: backing(DiscountBacking).Name의 프로모션이 활성인지 판정
package price

// HasActivePromotion checks if the promotion named by backing.Name is active.
func HasActivePromotion(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*PriceContext)
	db := backing.(*DiscountBacking)
	for _, p := range ctx.Promotions {
		if p.Name == db.Name && p.Active {
			return true, db
		}
	}
	return false, nil
}
